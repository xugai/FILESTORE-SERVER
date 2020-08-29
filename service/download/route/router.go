package route

import (
	"FILESTORE-SERVER/service/download/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")
	router.Use(cors.Default())  // 用于跨域访问
	router.POST("/file/downloadurl", client.DownloadURLHandler)

	return router
}
