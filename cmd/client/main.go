package main

import (
	"cache/rpc"
	"context"
	"fmt"
	"log"
)

type Args struct {
	Num1, Num2 int
}

func main() {
	client, err := rpc.Dial("tcp", ":7001")
	if err != nil {
		log.Fatalln("network dial error: ", err)
	}
	args := &Args{
		1,
		2,
	}
	var reply int
	err = client.Call(context.Background(), "Foo.Sum", args, &reply)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reply)
}
