package http

import (
	"cache/consistenthash"
	"cache/rpc"
	"fmt"
	"log"
	"sync"
)

const (
	defaultReplicas = 50
)

type RpcPool struct {
	self    string
	mu      sync.Mutex
	peers   *consistenthash.Map    // 一致性哈希
	clients map[string]*rpc.Client // 保存节点
}

var rpcPoll *RpcPool

func NewRpcPool(self string) *RpcPool {
	rpcPoll = &RpcPool{
		self:    self,
		peers:   consistenthash.New(defaultReplicas, nil),
		clients: make(map[string]*rpc.Client),
	}
	return rpcPoll
}

func GetRpcPool() *RpcPool {
	return rpcPoll
}

func (p *RpcPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// SetPeers 初始化consistent hash，并添加节点
func (p *RpcPool) SetPeers(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.clients = make(map[string]*rpc.Client, len(peers))
	for _, peer := range peers {
		p.clients[peer], _ = rpc.Dial("tcp", peer)
	}
}

// PickPeer 选择一个节点，实现PeerPicker接口
func (p *RpcPool) PickPeer(key string) (*rpc.Client, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" {
		p.Log("Pick rpc pool %s", peer)
		return p.clients[peer], true
	}

	return nil, false
}
