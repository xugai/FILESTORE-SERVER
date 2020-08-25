package route

import (
	"FILESTORE-SERVER/service/download/client"
	"github.com/gin-gonic/gin"
)



func Router() *gin.Engine {
	router := gin.Default()

	router.POST("/file/download", client.DownloadURLHandler)

	return router
}
