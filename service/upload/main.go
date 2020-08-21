package main

import (
	"FILESTORE-SERVER/service/upload/handler"
	"FILESTORE-SERVER/service/upload/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
)

func startupUploadService() {
	newRegistry := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service := micro.NewService(
		micro.Registry(newRegistry),
		micro.Name("go.micro.service.upload"),
	)
	service.Init()
	proto.RegisterUploadServiceHandler(service.Server(), new(handler.Upload))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startupUploadService()
}
