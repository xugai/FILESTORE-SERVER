package route

import (
	"FILESTORE-SERVER/service/upload/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Static("/static/", "./static")
	router.Use(cors.Default())

	router.POST("/file/update", client.UpdateFileMetaHandler)
	router.GET("/file/upload", client.UploadHandler)
	router.POST("/file/upload", client.UploadPostHandler)

	// 分块上传
	router.POST("/file/mpupload/init", client.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart", client.UploadChunkFileHandler)
	router.POST("/file/mpupload/cancel", client.CancelUploadHandler)
	router.POST("/file/mpupload/complete", client.CompleteMultipartUploadHandler)

	// 文件秒传
	router.POST("/file/fastupload", client.TryFastUploadHandler)
	return router
}
