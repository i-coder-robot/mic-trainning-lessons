package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/biz"
	"github.com/i-coder-robot/mic-trainning-lessons/account_srv/proto/pb"
	"github.com/i-coder-robot/mic-trainning-lessons/internal"
	"github.com/i-coder-robot/mic-trainning-lessons/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	//ip := flag.String("ip", "127.0.0.1", "输入ip")
	//port := flag.Int("port", 9095, "输入端口")
	//flag.Parse()
	//addr := fmt.Sprintf("%s:%d", *ip, *port)
	port := util.GenRandomPort()
	if internal.AppConf.Debug {
		port = internal.AppConf.AccountSrvConfig.Port
	}
	addr := fmt.Sprintf("%s:%d", internal.AppConf.AccountSrvConfig.Host, port)
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv启动异常:" + err.Error())
		panic(err)
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	checkAddr := fmt.Sprintf("%s:%d", internal.AppConf.AccountSrvConfig.Host, port)
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	randUUID := uuid.New().String()
	reg := api.AgentServiceRegistration{
		Name:    internal.AppConf.AccountSrvConfig.SrvName,
		ID:      randUUID,
		Port:    port,
		Tags:    internal.AppConf.AccountSrvConfig.Tags,
		Address: internal.AppConf.AccountSrvConfig.Host,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port))
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
