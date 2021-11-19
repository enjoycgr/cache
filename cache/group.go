package cache

import (
	"cache/lru"
	"cache/peer"
	"fmt"
	"log"
	"sync"
)

// 未命中获取数据源的回调
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	mainCache *cache // 本地的缓存
	getter    Getter // 未命中缓存时的回调
	peers     peer.PeerPicker
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

func (g *Group) RegisterPeers(peers peer.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

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

func (g *Group) Add(key string, value string) {
	g.mainCache.add(key, ByteView{[]byte(value)})
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 从远程获取
	if g.peers != nil {
		if p, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(p, key); err == nil {
				return value, nil
			}
			log.Println("[Cache] Failed to get from peer", err)
		}
	}

	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer peer.PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}

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
