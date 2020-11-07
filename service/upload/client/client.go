package client

import (
	"FILESTORE-SERVER/service/upload/handler"
	"FILESTORE-SERVER/service/upload/proto"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-plugins/registry/kubernetes/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var uploadCli proto.UploadService
var timeoutSetting client.CallOption = func(o *client.CallOptions) {
	o.RequestTimeout = 40 * time.Second
	o.DialTimeout = 40 * time.Second
}

func init() {
	k8sRegistry := kubernetes.NewRegistry()
	//newRegistry := consul.NewRegistry(registry.Addrs("192.168.10.3:8500"))
	newService := micro.NewService(
		micro.Registry(k8sRegistry),
	)
	newService.Init()
	c := newService.Client()
	c.Init(grpc.MaxSendMsgSize(10 * 1024 * 1024))
	uploadCli = proto.NewUploadService("go.micro.service.upload", c)
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

func UpdateFileMetaHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	if opType != "0" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg": "Request parameters invalid!",
		})
		return
	}
	userName := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")

	respUpdateFileMeta, err := uploadCli.UpdateFileMeta(context.TODO(), &proto.ReqUpdateFileMeta{
		Username:    userName,
		Filehash:    filehash,
		Newfilename: newFileName,
	})
	if err != nil || respUpdateFileMeta.Code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "File meta info update failed, please check log to get more details!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respUpdateFileMeta.Code,
		"msg": respUpdateFileMeta.Msg,
	})
}

func InitialMultipartUploadHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Invalid request parameter, please check log to get more details!",
		})
		log.Printf("Invalid request parameter: filesize, please check!")
		return
	}
	respInitialMultipartUpload, err := uploadCli.InitialMultipartUpload(context.TODO(), &proto.ReqInitialMultipartUpload{
		Username: username,
		Filehash: filehash,
		Filesize: int64(filesize),
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	multipartUploadInfo := handler.MultipartUploadInfo{}
	_ = json.Unmarshal(respInitialMultipartUpload.Initialresult, &multipartUploadInfo)
	c.JSON(http.StatusOK, gin.H{
		"code": respInitialMultipartUpload.Code,
		"msg": respInitialMultipartUpload.Msg,
		"data": multipartUploadInfo,
	})
}


//todo 如果不在断点执行的情况下，分块上传的请求会失败。这很奇怪......
//todo 除此之外，原本3MB左右的文件，通过上传4个分块后，发现4个分块的大小都是默认的分块大小1MB，也就是说，没有正确地读取文件该有的大小的字节
//todo 对于以上问题，下一步需要做的是：拿之前写的测试分块上传ws的测试代码，然后模拟现在web端的操作，尝试进行分块上传，验证是否会出现同样的问题，以此来诊断到底问题是出在前端，还是出在微服务化之后
func UploadChunkFileHandler(c *gin.Context) {
	uploadId := c.Request.FormValue("uploadid")
	chkIdx, err := strconv.Atoi(c.Request.FormValue("index"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Upload chunk file error, please check log to get more details",
		})
		log.Println("The request parameters are invalid")
		return
	}
	//buf := make([]byte, 1024 * 100)
	//todo 问题发现：好像前端传到客户端这边的字节虽然是符合预期的，但是在读取到另一个字节数组中时，并不是一次read就能够读取完的。读取到一定数量的字节后会自己退出方法？？？
	all, _ := ioutil.ReadAll(c.Request.Body)
	//n, err := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload chunk file error, please check log to get more details",
		})
		log.Printf("Read from request body error: %v\n", err)
		return
	}
	respUploadChunkFile, err := uploadCli.UploadChunkFile(context.TODO(), &proto.ReqUploadChunkFile{
		Uploadid:   uploadId,
		Chkidx:     int32(chkIdx),
		Chkcontent: all,
	}, timeoutSetting)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respUploadChunkFile.Code,
		"msg": respUploadChunkFile.Msg,
	})
}

func CompleteMultipartUploadHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	uploadId := c.Request.FormValue("uploadid")
	fileHash := c.Request.FormValue("filehash")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Upload chunk file error, please check log to get more details",
		})
		log.Println("The request parameters are invalid")
		return
	}
	fileName := c.Request.FormValue("filename")
	respCompleteMultipartUpload, err := uploadCli.CompleteMultipartUpload(context.TODO(), &proto.ReqCompleteMultipartUpload{
		Username: userName,
		Uploadid: uploadId,
		Filehash: fileHash,
		Filesize: int64(fileSize),
		Filename: fileName,
	}, timeoutSetting)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respCompleteMultipartUpload.Code,
		"msg": respCompleteMultipartUpload.Msg,
	})
}

func CancelUploadHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	respCancelUpload, err := uploadCli.CancelUpload(context.TODO(), &proto.ReqCancelUpload{
		Filehash: fileHash,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respCancelUpload.Code,
		"msg": respCancelUpload.Msg,
	})
}

func TryFastUploadHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Fast upload file error, please check log to get more details!",
		})
		log.Println(err)
		return
	}
	respFastUpload, err := uploadCli.FastUpload(context.TODO(), &proto.ReqFastUpload{
		Username: userName,
		Filehash: fileHash,
		Filename: fileName,
		Filesize: int64(fileSize),
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Fast upload file error, please check log to get more details!",
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respFastUpload.Code,
		"msg": respFastUpload.Msg,
	})
}