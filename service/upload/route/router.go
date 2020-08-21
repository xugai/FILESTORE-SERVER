package route

import (
	"FILESTORE-SERVER/service/upload/client"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.GET("/file/upload", client.UploadHandler)
	router.POST("/file/upload", client.UploadPostHandler)
	return router
}
