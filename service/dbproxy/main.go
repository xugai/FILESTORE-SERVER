package main

import (
	"FILESTORE-SERVER/service/dbproxy/handler"
	"FILESTORE-SERVER/service/dbproxy/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
	"time"
)

func startupDBService() {

	k8sRegistry := kubernetes.NewRegistry()

	//newRegistry := consul.NewRegistry(registry.Addrs("192.168.10.3:8500"))
	service := micro.NewService(
		micro.Registry(k8sRegistry),
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
