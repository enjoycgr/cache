package peer

import (
	"cache/cache"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_cache/"

type HttpPool struct {
	self     string
	basePath string
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

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path:" + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	if r.Method == "GET" {
		p.get(w, r)
	}

	if r.Method == "POST" {
		p.post(w, r)
	}
}

// 获取缓存
func (p *HttpPool) get(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	groupName := parts[0]
	key := parts[1]

	group := cache.GetGroup(groupName)
	if group == nil {
		http.Error(w, strings.Join([]string{"no such group", groupName}, ": "), http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// 设置缓存
func (p *HttpPool) post(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 1)
	if len(parts) != 1 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	groupName := parts[0]

	group := cache.GetGroup(groupName)
	if group == nil {
		http.Error(w, strings.Join([]string{"no such group", groupName}, ": "), http.StatusNotFound)
		return
	}
	r.ParseForm()
	key := r.PostForm.Get("key")
	value := r.PostForm.Get("value")
	log.Println(key, value)
	group.Add(key, value)
	w.WriteHeader(http.StatusOK)
}
