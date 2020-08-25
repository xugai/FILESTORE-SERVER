package auth

import (
	dbCli "FILESTORE-SERVER/service/dbproxy/client"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ifTokenIsValid(username string, token string) bool {
	// expired time: 30 minutes
	half_an_hour := 60 * 30
	parseUint, err2 := strconv.ParseUint(token[(len(token)-8):], 16, 32)
	if err2 != nil {
		fmt.Printf("Convert from hex to dec failed: %v\n", err2)
		return false
	}
	execResult, err2 := dbCli.IfTokenIsValid(username, token)
	if err2 != nil {
		log.Println(err2)
		return false
	}
	return time.Now().Unix() - int64(parseUint) < int64(half_an_hour) && execResult.Suc
}

func IdentityMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		userName := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		if len(userName) < 3 || !ifTokenIsValid(userName, token) {
			c.JSON(http.StatusForbidden, nil)
			c.Abort()
		}
		c.Next()
	}
}
