package handler

import (
	"FILESTORE-SERVER/db"
	userProto "FILESTORE-SERVER/service/account/proto"
	downloadProto "FILESTORE-SERVER/service/download/proto"
	uploadProto "FILESTORE-SERVER/service/upload/proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"net/http"
)

var (
	userCli userProto.UserService
	uploadCli uploadProto.UploadService
	downloadCli downloadProto.DownloadService
)

func init() {
	newRegistry := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Registry(newRegistry),
		)
	// 初始化，解析命令行参数端
	service.Init()

	// 初始化一个account service的客户端
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
	uploadCli = uploadProto.NewUploadService("go.micro.service.upload", service.Client())
	downloadCli = downloadProto.NewDownloadService("go.micro.service.download", service.Client())
}

func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func SignupPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	passWord := c.Request.FormValue("password")
	respSignup, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: userName,
		Password: passWord,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respSignup.Code,
		"msg": respSignup.Message,
	})
}

func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func SigninPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	passWord := c.Request.FormValue("password")
	respSignin, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: userName,
		Password: passWord,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	uploadEntry, err := getUploadEntry()
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respSignin.Code,
		"data": struct {
			Token string
			Username string
			Location string
			UploadEntry string
		}{
			Token: respSignin.Token,
			Username: userName,
			Location: "/static/view/home.html",
			UploadEntry: uploadEntry,
		},
		"msg": respSignin.Message,
	})
}

func UserInfoPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	respUserInfo, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: userName,
	})
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": respUserInfo.Code,
		"msg": respUserInfo.Message,
		"data": db.User{
			UserName: respUserInfo.Username,
			Email: respUserInfo.Email,
			Phone: respUserInfo.Phone,
			LastActive: respUserInfo.LastActiveAt,
			SignupAt: respUserInfo.Signup,
			Status: int(respUserInfo.Status),
		},
	})
}

func getUploadEntry() (string, error) {
	respUploadEntry, err := uploadCli.UploadEntry(context.TODO(), &uploadProto.ReqUploadEntry{})
	if err != nil {
		return "", err
	}
	return respUploadEntry.Entry, nil
}
