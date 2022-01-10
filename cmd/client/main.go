package main

import (
	"cache/rpc"
	"cache/server"
	"context"
	"fmt"
	"log"
	"sync"
)

func main() {
	client, err := rpc.Dial("tcp", ":8001")
	if err != nil {
		log.Fatalln("network dial error: ", err)
	}
	args := server.SetRequest{
		Key:   "1",
		Value: "test",
	}
	var reply server.SetResponse
	var i int
	var wg sync.WaitGroup
	for {
		i++
		wg.Add(1)
		go func() {
			err = client.Call(context.Background(), "Server.Set", args, &reply)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(reply)
			wg.Done()
		}()

		if i > 10 {
			break
		}
	}
	wg.Wait()
}
