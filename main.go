package main

import (
	"cache/cache"
	"cache/peer"
	"fmt"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	cache.NewGroup("scores", 2<<20, cache.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}),
	)

	addr := ":9999"
	peers := peer.NewHttpPool(addr)
	http.ListenAndServe(addr, peers)
}
