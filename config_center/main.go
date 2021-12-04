package main

import (
	"fmt"
	"github.com/i-coder-robot/mic-trainning-lessons/internal"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	nacosConfig := internal.ViperConf.NacosConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port:   nacosConfig.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		//NamespaceId: "3a6aa0c4-b492-4624-bb59-16cc3f27416f",
		NamespaceId:         nacosConfig.NameSpace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	//content,err:=configClient.GetConfig(vo.ConfigParam{
	//	DataId:"account_srv.json",
	//	Group: "dev",
	//})

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
