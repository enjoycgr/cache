package cache

import (
	"cache/cachepb"
	"cache/consistenthash"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	defaultBasePath = "/_cache/"
	defaultReplicas = 50
)

type HttpPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map    // 一致性哈希
	httpGetters map[string]*httpGetter // 保存节点
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// Set 初始化consistent hash，并添加节点
func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 选择一个节点，实现PeerPicker接口
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick httppool %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _ PeerPicker = (*HttpPool)(nil)

type httpGetter struct {
	baseURL string
}

// Get 从其他节点获取缓存
func (h *httpGetter) Get(request *cachepb.GetRequest, response *cachepb.GetResponse) error {
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(request.GetGroup()),
		url.QueryEscape(request.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reding response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, response); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

// Set 从其他节点中设置缓存
func (h *httpGetter) Set(request *cachepb.SetRequest, response *cachepb.SetResponse) error {
	u := fmt.Sprintf("%v%v",
		h.baseURL,
		url.QueryEscape(request.GetGroup()),
	)
	res, err := http.PostForm(u, url.Values{
		"key":   {request.GetKey()},
		"value": {request.GetValue()},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	response.Res = true

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)
