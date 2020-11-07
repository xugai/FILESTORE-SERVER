package route

import (
	"FILESTORE-SERVER/asset"
	securityLayer "FILESTORE-SERVER/service/apigw/auth"
	"FILESTORE-SERVER/service/apigw/handler"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset: asset.Asset,
		AssetDir: asset.AssetDir,
		AssetInfo: asset.AssetInfo,
		Prefix: root,
	}
	return &binaryFileSystem{
		fs: fs,
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	// todo 也就是说，通过ide开启路由跟通过命令行开启路由，前者因为相对路径的缘故因此可以索引到静态资源文件；后者此时就无法索引到了
	//router.Static("/static/", "./static")
	// 将前端静态文件打包到bin文件里边
	router.Use(static.Serve("/static/", BinaryFileSystem("static")))
	// 使用gin插件支持跨域请求
	//router.Use(cors.Default())

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
