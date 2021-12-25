package main

import (
	"cache/rpc"
	"log"
	"net"
)

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func main() {
	var foo Foo
	rpc.Register(foo)

	l, err := net.Listen("tcp", ":7001")
	if err != nil {
		log.Fatal("network error: ", err)
	}

	rpc.Accept(l)
}
