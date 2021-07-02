package main

import (
	"context"
	"ithaiq/grpc/client"
	"ithaiq/grpc/server"
	"log"
	"net"
	"sync"
	"time"
)

type Grpc int

type GrpcReq struct {
	Num1, Num2 int
}

func (f Grpc) Test(args GrpcReq, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Grpc
	if err := server.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	server.Accept(l)
}

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
			args := &GrpcReq{Num1: i, Num2: i * i}
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			var reply int
			if err := conn.Call(ctx, "Grpc.Test", args, &reply); err != nil {
				log.Fatal("call Grpc.Test error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
