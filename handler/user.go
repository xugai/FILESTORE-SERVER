package handler

import (
	"FILESTORE-SERVER/config"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// for gin web fwk
func SignupPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	passWord := c.Request.FormValue("password")
	// validation
	if len(userName) < 3 || len(passWord) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Invalid request parameter(s)!",
		})
		return
	}
	// encrypt for password
	result := db.UserSignup(userName, utils.Sha1([]byte(passWord + config.Salt)))
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Sign up failed!",
		})
	}
}

func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// for gin web fwk
func SigninPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	passWord := c.Request.FormValue("password")

	result := db.UserSignin(userName, utils.Sha1([]byte(passWord + config.Salt)))
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Login failed, your username or password maybe is incorrect!",
		})
		return
	}
	token := utils.GenToken(userName)
	result = db.FlushUserToken(userName, token)
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Flush user token failed, please check log",
		})
		return
	}
	serverResponse := utils.NewServerResponse(200, "OK", struct {
		Location string
		Username string
		Token    string
	}{
		Location: "/static/view/home.html",
		Username: userName,
		Token:    token,
	})
	c.Data(http.StatusOK, "application/json", serverResponse.GetInByteStream())
}

// for gin web fwk
func UserInfoPostHandler(c *gin.Context) {
	userName := c.Request.FormValue("username")
	userInfo, err := db.GetUserInfo(userName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Get user info failed, please check!",
		})
		return
	}
	serverResponse := utils.NewServerResponse(200, "OK", userInfo)
	c.Data(http.StatusOK, "application/json", serverResponse.GetInByteStream())
}
