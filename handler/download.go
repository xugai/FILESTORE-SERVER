package handler

import (
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/store/oss"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	objectKeyPrefix = "oss/image/"
)

func DownloadURLHandler(c *gin.Context) {
	//1. 获得用户传过来的filehash
	filehash := c.Request.FormValue("filehash")
	//2. 通过filehash去db里面查询相应的文件的path
	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Get download url failed, please check log to get more details!",
		})
		log.Printf("%v\n", err)
		return
	}
	//3. 然后通过filepath获得ali oss的signed download url
	signedURL := oss.Download(objectKeyPrefix + fileMeta.FileName.String)
	//4. 最后将signed url返回给用户
	c.JSON(http.StatusOK, signedURL)
}
