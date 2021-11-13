package main

import (
	"flag"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/biz"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/internal"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"google.golang.org/grpc"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入ip")
	port := flag.Int("port", 9095, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
