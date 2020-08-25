package client

import (
	"FILESTORE-SERVER/service/download/proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"net/http"
)

var downloadCli proto.DownloadService

func init() {
	newRegistry := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service := micro.NewService(
		micro.Registry(newRegistry),
	)
	service.Init()
	downloadCli = proto.NewDownloadService("go.micro.service.download", service.Client())
}

func DownloadURLHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	respDownloadURL, err := downloadCli.DownloadURL(context.TODO(), &proto.ReqDownloadURL{
		Filehash: fileHash,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": respDownloadURL.Code,
			"msg": respDownloadURL.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, respDownloadURL.Url)
}
