package main

import (
	"cache/rpc"
	"cache/server"
	"context"
	"fmt"
	"log"
)

func main() {
	client, err := rpc.Dial("tcp", ":7001")
	if err != nil {
		log.Fatalln("network dial error: ", err)
	}
	args := server.GetRequest{
		Key: "test",
	}
	var reply server.GetResponse
	err = client.Call(context.Background(), "Service.Get", args, &reply)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply)
}
