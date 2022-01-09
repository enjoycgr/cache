package main

import (
	"cache/app"
	"cache/rpc"
	"cache/server"
	"log"
	"net"
)

type Server struct {
}

func (s Server) Get(request server.GetRequest, response *server.GetResponse) error {
	value, _ := cache.Get(request.Key)
	response.Value = value.String()
	return nil
}

func (s Server) Set(request server.SetRequest, response *server.SetResponse) error {
	cache.Set(request.Key, app.ByteView{B: []byte(request.Value)})
	response.Result = true
	return nil
}

var cache *app.Cache

func main() {
	cache = app.NewCache(2<<20, nil)

	var service Server
	//var service Foo
	rpc.Register(service)

	l, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal("network error: ", err)
	}

	rpc.Accept(l)
}
