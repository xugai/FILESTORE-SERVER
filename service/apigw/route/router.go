package route

import (
	securityLayer "FILESTORE-SERVER/service/apigw/auth"
	"FILESTORE-SERVER/service/apigw/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	// todo 也就是说，通过ide开启路由跟通过命令行开启路由，前者因为相对路径的缘故因此可以索引到静态资源文件；后者此时就无法索引到了
	router.Static("/static/", "./static")
	// 使用gin插件支持跨域请求
	router.Use(cors.Default())

	router.GET("/user/signup", handler.SignupHandler)
	router.POST("/user/signup", handler.SignupPostHandler)
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.SigninPostHandler)

	// 使用鉴权中间件
	router.Use(securityLayer.IdentityMiddleware())
	router.POST("/user/info", handler.UserInfoPostHandler)
	router.POST("/file/query", handler.QueryFileMetasHandler)
	return router
}
