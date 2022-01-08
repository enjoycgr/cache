package main

import (
	"cache/core/cache"
	"cache/rpc"
	"log"
	"net"
)

type Foo int

type Args1 struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args1, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func main() {
	var group cache.Group
	rpc.Register(group)

	l, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Fatal("network error: ", err)
	}

	rpc.Accept(l)
}
