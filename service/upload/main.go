package main

import (
	"FILESTORE-SERVER/service/upload/config"
	"FILESTORE-SERVER/service/upload/handler"
	"FILESTORE-SERVER/service/upload/proto"
	"FILESTORE-SERVER/service/upload/route"
	"github.com/micro/go-micro/v2"
	//"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server/grpc"
	//"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
)

func startupUploadService() {

	k8sRegistry := kubernetes.NewRegistry()

	//newRegistry := consul.NewRegistry(registry.Addrs("192.168.10.3:8500"))
	service := micro.NewService(
		micro.Registry(k8sRegistry),
		micro.Name("go.micro.service.upload"),
	)
	service.Init()
	s := service.Server()
	s.Init(grpc.MaxMsgSize(10 * 1024 * 1024))
	proto.RegisterUploadServiceHandler(s, new(handler.Upload))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func startupUploadServiceClient() {
	router := route.Router()
	router.Run(config.UploadServiceHost)
}

func main() {
	go startupUploadServiceClient() // 开启upload service rpc client
	startupUploadService() // 开启upload service rpc server
}
