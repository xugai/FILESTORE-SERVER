package main

import (
	"FILESTORE-SERVER/service/account/handler"
	"FILESTORE-SERVER/service/account/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
)

func main() {
	// 创建一个service
	newRegistry := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
		)
	service := micro.NewService(
		micro.Registry(newRegistry),
		micro.Name("go.micro.service.user"),
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Printf("%v\n", err)
	}

}
