package main

import (
	"FILESTORE-SERVER/service/dbproxy/handler"
	"FILESTORE-SERVER/service/dbproxy/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"time"
)

func startupDBService() {
	newRegistry := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service := micro.NewService(
		micro.Registry(newRegistry),
		micro.Name("go.micro.service.dbproxy"),
		micro.RegisterTTL(10*time.Second), // 声明超时时间，避免consul没有主动删除已失去心跳的节点
		micro.RegisterInterval(5*time.Second),
	)
	service.Init()

	proto.RegisterDBProxyServiceHandler(service.Server(), new(handler.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startupDBService()
}
