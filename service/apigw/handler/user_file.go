package handler

import (
	"FILESTORE-SERVER/service/account/proto"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func QueryFileMetasHandler(c *gin.Context) {
	limitCnt, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg": "Request parameters are invalid!",
		})
		return
	}
	userName := c.Request.FormValue("username")
	respQueryFileMetas, err := userCli.QueryFileMetas(context.TODO(), &proto.ReqQueryFileMetas{
		Username: userName,
		LimitCnt: int32(limitCnt),
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Failed",
		})
		return
	}
	c.Data(http.StatusOK, "application/json", respQueryFileMetas.FileMetas)
}
