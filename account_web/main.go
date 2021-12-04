package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i-coder-robot/mic-trainning-lessons/account_web/handler"
	"github.com/i-coder-robot/mic-trainning-lessons/internal"
)

func init() {
	err := internal.Reg(internal.AppConf.AccountWebConfig.Host,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.Port,
		internal.AppConf.AccountWebConfig.Tags)
	if err != nil {
		panic(err)
	}
}

func main() {
	ip := flag.String("ip", "192.168.0.106", "输入Ip")
	port := flag.Int("port", 8081, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	r := gin.Default()
	accountGroup := r.Group("/v1/account")
	{
		accountGroup.GET("/list", handler.AccountListHandler)
		accountGroup.POST("/login", handler.LoginByPasswordHandler)
		accountGroup.GET("/captcha", handler.CaptchaHandler)
	}
	r.GET("/health", handler.HealthHandler)

	r.Run(addr)
}
