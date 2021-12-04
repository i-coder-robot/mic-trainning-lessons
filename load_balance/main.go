package main

import (
	"context"
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	//addr := fmt.Sprintf("%s:%d", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port)
	addr := fmt.Sprintf("127.0.0.1:8500")
	dialAddr := fmt.Sprintf("consul://%s/account_srv?wait=14", addr)
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Fatal(err)
	}

	defer conn.Close()
	client := pb.NewAccountServiceClient(conn)
	res, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		panic(err)
	}
	for idx, item := range res.AccountList {
		fmt.Println(fmt.Sprintf("%d---%v", idx, item))
	}
}
