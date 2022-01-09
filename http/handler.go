package http

import (
	"cache/server"
	"context"
	"net/http"
)

// Get 获取缓存
func Get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	if key == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	client, ok := rpcPoll.PickPeer(key)
	if ok != true {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	args := server.GetRequest{Key: key}
	reply := server.GetResponse{}
	err := client.Call(context.Background(), "Server.Get", args, &reply)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reply.Value == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(reply.Value))
}

// Set 设置缓存
func Set(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	value := query.Get("value")
	if key == "" || value == "" {
		http.Error(w, "key or value is null", http.StatusInternalServerError)
		return
	}

	client, ok := rpcPoll.PickPeer(key)
	if ok != true {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	args := server.SetRequest{Key: key, Value: value}
	reply := server.GetResponse{}
	err := client.Call(context.Background(), "Server.Set", args, &reply)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(reply.Value))
}
