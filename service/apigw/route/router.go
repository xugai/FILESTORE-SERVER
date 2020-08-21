package route

import (
	"FILESTORE-SERVER/service/apigw/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")
	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.SignupPostHandler)
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.SigninPostHandler)
	router.POST("/user/info", handler.UserInfoPostHandler)

	return router
}
