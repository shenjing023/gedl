package main

import (
	"context"
	"fmt"
	"gedl/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type HelloService struct {
	pb.UnimplementedHelloServer
}

func (a HelloService) HelloWorld(ctx context.Context, request *pb.Request) (*pb.Reply, error) {
	fmt.Println("svc1")
	return &pb.Reply{Output: "svc1"}, nil
	// fmt.Println("svc2")
	// return &pb.Reply{Output: "svc2"}, nil
}

const (
	host = "localhost"
	port = "5003"
	Addr = "0.0.0.0:5003"
)

func main() {
	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	var srv = &HelloService{}
	pb.RegisterHelloServer(grpcServer, srv)
	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	conf := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: time.Second * 5,
	}
	r, err := NewRegister(conf, "svc", "hello", host, port)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	fmt.Println("启动成功 === > ", Addr)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for a := range c {
		switch a {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("退出")
			fmt.Println(r.Close())
			return
		default:
			return
		}
	}

}
