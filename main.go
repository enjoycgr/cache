package main

import (
	"cache/cache"
	_ "cache/config"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<20, cache.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}),
	)
}

func startCacheServer(addr string, addrs []string, c *cache.Group) {
	peers := cache.NewHttpPool(addr)
	// 添加其他节点到哈希环
	peers.Set(addrs...)
	c.RegisterPeers(peers)
	log.Println("cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, c *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := c.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("font-end server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	server := viper.Get("server").([]interface{})
	apihost := viper.Get("apihost").(string)
	for _, v := range server {
		fmt.Println(v)
	}
	peers := cache.NewHttpPool(apihost)
	err := http.ListenAndServe(apihost, peers)
	if err != nil {
		log.Println("api server start err:", err)
	}
}
