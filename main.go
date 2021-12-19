package main

import (
	"cache/cache"
	_ "cache/config"
	"cache/rpc"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
	"1":    "1",
	"2":    "2",
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

// addr：本地节点，addrs：其他节点切片，c：group
// 开启缓存服务
func startCacheServer(addr string, addrs []string, c *cache.Group) {
	peers := cache.NewHttpPool(addr)
	// 添加其他节点到哈希环
	peers.Set(addrs...)
	c.RegisterPeers(peers)
	log.Printf("cache is running at %s \n", addr)
	log.Fatal(http.ListenAndServe(":8001", peers))
}

// 开启api服务，和用户交互
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
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func startRpcServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	rpc.Accept(l)
}

func main() {
	//server := viper.GetStringSlice("server")
	//apihost := viper.GetString("apihost")
	//cachehost := viper.GetString("cachehost")
	//
	//group := createGroup()
	//go startCacheServer(cachehost, server, group)
	//
	//startAPIServer(apihost, group)
	log.SetFlags(0)
	addr := make(chan string)
	go startRpcServer(addr)
	client, _ := rpc.Dial("tcp", <-addr)
	defer func() {
		_ = client.Close()
	}()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("rpc req %d", i)
			var reply string
			//if err := client.Go("Foo.Sum", args, &reply, nil); err != nil {
			//	log.Fatal("call Foo.Sum error:", err)
			//}
			Call := client.Go("Foo.Sum", args, &reply, nil)
			<-Call.Done
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()

}
