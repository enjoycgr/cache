package main

import (
	"cache/server"
)

func main() {
	server.StartRpcServer(":8001")
}
