package cache

import (
	"cache/core/cachepb"
	"cache/lru"
	"fmt"
	"log"
	"sync"
)

// Getter 未命中获取数据源的回调
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 接口型函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	mainCache *cache // 本地的缓存
	getter    Getter // 未命中缓存时的回调
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	groups[name] = &Group{
		name: name,
		mainCache: &cache{
			cacheBytes: cacheBytes,
			lru:        lru.New(cacheBytes, nil),
		},
		getter: getter,
	}

	return groups[name]
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// Get 从本地获取缓存，获取不到从其他节点load缓存
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("Cache hit")
		return v, nil
	}

	return g.load(key)
}

// Set 设置缓存
func (g *Group) Set(key string, value string) error {
	if g.peers != nil {
		if p, ok := g.peers.PickPeer(key); ok {
			return g.setFromPeer(p, key, value)
		}
	}
	view := ByteView{b: cloneBytes([]byte(value))}
	g.populateCache(key, view)
	return nil
}

// 从其他节点获取缓存，获取不到调用回调方法
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if p, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(p, key); err == nil {
				return value, nil
			}
			log.Println("[Cache] Failed to get from httppool", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &cachepb.GetRequest{
		Group: g.name,
		Key:   key,
	}
	res := &cachepb.GetResponse{}
	err := peer.Get(req, res)

	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: res.Value}, nil
}

func (g *Group) setFromPeer(peer PeerGetter, key string, value string) error {
	req := &cachepb.SetRequest{
		Group: g.name,
		Key:   key,
		Value: value,
	}
	res := &cachepb.SetResponse{}
	err := peer.Set(req, res)
	return err
}

// 回调函数获取value，并存入缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
