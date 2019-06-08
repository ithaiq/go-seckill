package main

import (
	"ithaiq/gtcp/gnet"
)

func main() {
	server := gnet.NewServer("gtcp")
	server.Serve()
}
