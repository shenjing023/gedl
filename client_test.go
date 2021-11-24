package main

import (
	"context"
	"fmt"
	"gedl/pb"
	"log"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func TestClient(t *testing.T) {
	conf := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	}
	r, err := NewDiscovery(conf, "svc", "hello")
	if err != nil {
		panic(err)
	}
	resolver.Register(r)
	// 连接服务器
	conn, err := grpc.DialContext(context.Background(),
		GetPrefix("svc", "hello"),
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelloClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 10; i++ {
		res, err := client.HelloWorld(ctx, &pb.Request{Input: "1111"})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.Output)
	}

}
