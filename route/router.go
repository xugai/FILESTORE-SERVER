package route

import (
	"FILESTORE-SERVER/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// 获得包含像Logger Recovery等的中间件的Engine
	e := gin.Default()

	// 处理静态资源
	e.Static("/static", "./static")

	// 处理不需要拦截器验证的请求
	// 处理用户模块的请求
	e.GET("/user/signup", handler.SignupHandler)
	e.POST("/user/signup", handler.SignupPostHandler)
	e.GET("/user/signin", handler.SigninHandler)
	e.POST("/user/signin", handler.SigninPostHandler)
	e.POST("/user/info", handler.UserInfoPostHandler)

	// 加入中间件，用来进行身份验证
	e.Use(handler.IdentityMiddleware())

	// 文件相关接口
	e.GET("/file/upload", handler.UploadHandler)
	e.POST("/file/upload", handler.UploadPostHandler)
	e.GET("/file/upload/suc", handler.UploadFileSucHandler)
	e.GET("/file/metainfo", handler.GetFileMetaHandler)
	e.POST("/file/query", handler.QueryFileMetasHandler)
	e.POST("/file/download", handler.FileDownloadHandler)
	e.POST("/file/update", handler.FileMetaUpdateHandler)
	e.POST("/file/delete", handler.FileMetaDeleteHandler)
	e.POST("/file/fastupload", handler.TryFastUploadHandler)
	e.POST("/file/downloadurl", handler.DownloadURLHandler)
	e.POST("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	e.POST("/file/mpupload/mppart", handler.UploadChunkFileHandler)
	e.POST("/file/mpupload/complete", handler.CompleteUploadHandler)
	e.POST("/file/mpupload/cancel", handler.CancelUploadHandler)
	return e
}
