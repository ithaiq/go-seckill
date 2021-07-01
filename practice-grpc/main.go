package main

import (
	"fmt"
	"ithaiq/grpc/client"
	"ithaiq/grpc/server"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

/*func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(server.DefaultOption) //首先发送 Option 进行协议交换
	cc := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Grpc.Test",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("grpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}*/

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	conn, _ := client.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("grpc req %d", i)
			var reply string
			if err := conn.Call("Grpc.Test", args, &reply); err != nil {
				log.Fatal("call Grpc.Test error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}