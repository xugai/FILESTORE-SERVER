package client

import (
	"FILESTORE-SERVER/service/upload/proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"io/ioutil"
	"log"
	"net/http"
)

var uploadCli proto.UploadService


func init() {
	newRegistry := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	newService := micro.NewService(
		micro.Registry(newRegistry),
	)
	newService.Init()
	uploadCli = proto.NewUploadService("go.micro.service.upload", newService.Client())
}

func UploadHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/index.html")
}

func UploadPostHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	userName := c.Request.FormValue("username")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Failed to get uploaded file data, please check log to get more details!",
		})
		log.Printf("Failed to get uploaded file data, err: %v\n", err)
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "The file will be uploaded is invalid!",
		})
		return
	}
	respUploadFile, err := uploadCli.UploadFile(context.TODO(), &proto.ReqUploadFile{
		Username:    userName,
		Filename:    header.Filename,
		Filecontent: bytes,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respUploadFile.Code,
		"msg": respUploadFile.Message,
	})
}

