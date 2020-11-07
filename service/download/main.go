package main

import (
	"FILESTORE-SERVER/service/download/config"
	"FILESTORE-SERVER/service/download/handler"
	"FILESTORE-SERVER/service/download/proto"
	"FILESTORE-SERVER/service/download/route"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
)

func startupDownloadService() {
	k8sRegistry := kubernetes.NewRegistry()
	//newRegistry := consul.NewRegistry(registry.Addrs("192.168.10.3:8500"))
	service := micro.NewService(
		micro.Registry(k8sRegistry),
		micro.Name("go.micro.service.download"),
	)
	service.Init()
	proto.RegisterDownloadServiceHandler(service.Server(), new(handler.Download))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func startupDownloadServiceClient() {
	router := route.Router()
	router.Run(config.DownloadServiceHost)
}

func main() {
	go startupDownloadServiceClient() // 开启 download service rpc client
	startupDownloadService() // 开启 download service rpc server
}
