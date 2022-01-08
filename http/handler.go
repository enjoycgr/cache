package http

import (
	"cache/core/cache"
	"cache/core/cachepb"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
)

func (p *http2.HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (p *http2.HttpPool) get(w http.ResponseWriter, r *http.Request) {
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

	body, err := proto.Marshal(&cachepb.GetResponse{
		Value: view.ByteSlice(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// 设置缓存
func (p *http2.HttpPool) post(w http.ResponseWriter, r *http.Request) {
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
	if err := group.Set(key, value); err != nil {
		http.Error(w, strings.Join([]string{"Set cache error", err.Error()}, ": "), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
