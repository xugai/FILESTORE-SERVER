package route

import (
	"FILESTORE-SERVER/service/download/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static/", "./static")
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
	}))
	router.POST("/file/download", client.DownloadFileHandler)
	router.POST("/file/downloadurl", client.DownloadURLHandler)
	router.GET("/file/download/range", client.RangeDownloadHandler)
	return router
}
