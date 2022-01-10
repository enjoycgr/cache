package server

import (
	"cache/app"
	"cache/rpc"
	"log"
	"net"
)

type SetRequest struct {
	Key   string
	Value string
}

type SetResponse struct {
	Result bool
}

type GetRequest struct {
	Key string
}

type GetResponse struct {
	Value string
}

type Server struct {
}

var cache *app.Cache

func (s Server) Get(request GetRequest, response *GetResponse) error {
	value, _ := cache.Get(request.Key)
	response.Value = value.String()
	return nil
}

func (s Server) Set(request SetRequest, response *SetResponse) error {
	cache.Set(request.Key, app.ByteView{B: []byte(request.Value)})
	response.Result = true
	log.Println("call set success!")
	return nil
}

func StartRpcServer(addr string) {
	cache = app.NewCache(2<<20, nil)

	var service Server
	rpc.Register(service)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("network error: ", err)
	}

	rpc.Accept(l)
}
