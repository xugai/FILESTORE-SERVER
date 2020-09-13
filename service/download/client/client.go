package client

import (
	dbCli "FILESTORE-SERVER/service/dbproxy/client"
	"FILESTORE-SERVER/service/download/config"
	"FILESTORE-SERVER/service/download/proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"net/http"
	"os"
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

func DownloadFileHandler(c *gin.Context) {
	fileName := c.Request.FormValue("filename")
	fileHash := c.Request.FormValue("filehash")
	respDownloadFile, err := downloadCli.DownloadFile(context.TODO(), &proto.ReqDownloadFile{
		Filehash: fileHash,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": respDownloadFile.Code,
			"msg": respDownloadFile.Msg,
		})
		return
	}
	c.Header("content-disposition", "attachment; filename=\"" + fileName + "\"")
	c.Data(http.StatusOK, "application/octect-stream", respDownloadFile.FileContent)
}

func RangeDownloadHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	userName := c.Request.FormValue("username")
	getFileMetaExecResult, getFileMetaErr := dbCli.GetFileMeta(fileHash)
	getUserFileMetaExecResult, getUserFileMetaErr := dbCli.GetUserFileMeta(userName, fileHash)
	if getFileMetaErr != nil || getUserFileMetaErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Server error",
		})
		return
	}
	fileMeta := dbCli.ToTableFile(getFileMetaExecResult.Data)
	userFile := dbCli.ToTableUserFile(getUserFileMetaExecResult.Data)
	fpath := config.TmpStoreDir + fileMeta.FileName.String
	log.Println("range-download-file-path: ", fpath)
	file, err := os.Open(fpath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Server error",
		})
		log.Println(err)
		return
	}
	defer file.Close()
	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	c.Writer.Header().Set("content-disposition", "attachment; filename=\"" + userFile.FileName + "\"")
	c.File(fpath)
}