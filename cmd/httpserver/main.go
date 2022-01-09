package main

import (
	_ "cache/config"
	http2 "cache/http"
	"cache/server"
	"log"
	"net/http"
	"time"
)

func createRpcPoll(addr string, addrs []string) *http2.RpcPool {
	peers := http2.NewRpcPool(addr)
	peers.SetPeers(addrs...)
	return peers
}

// 开启api服务，和用户交互
func startAPIServer(apiAddr string) {
	http.HandleFunc("/get", http2.Get)
	http.HandleFunc("/set", http2.Set)
	log.Println("font-end server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))
}

func main() {
	go server.StartRpcServer("127.0.0.1:8001")
	time.Sleep(1 * time.Second)
	_ = createRpcPoll("127.0.0.1:8001", []string{"127.0.0.1:8001"})

	startAPIServer("127.0.0.1:8002")
}
