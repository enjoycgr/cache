package main

import (
	"cache/rpc"
	"cache/server"
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	client, err := rpc.Dial("tcp", ":8001")
	if err != nil {
		log.Fatalln("network dial error: ", err)
	}
	args := server.SetRequest{
		Key:   strconv.Itoa(1),
		Value: strconv.Itoa(1),
	}
	var reply server.SetResponse
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Call(context.Background(), "Server.Set", args, &reply)
		fmt.Println(reply)
	}
}
