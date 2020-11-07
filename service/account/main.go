package main

import (
	"FILESTORE-SERVER/service/account/handler"
	"FILESTORE-SERVER/service/account/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
)

func main() {
	// 创建一个service
	k8sRegistry := kubernetes.NewRegistry()
	//newRegistry := consul.NewRegistry(
	//	registry.Addrs("192.168.10.3:8500"),
	//	)
	service := micro.NewService(
		micro.Registry(k8sRegistry),
		micro.Name("go.micro.service.user"),
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Printf("%v\n", err)
	}

}
