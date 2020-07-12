package handler

// 用户鉴权接口

import (
	"FILESTORE-SERVER/db"
	"fmt"
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
	return time.Now().Unix() - int64(parseUint) < int64(half_an_hour) && db.IfTokenIsValid(username, token)
}

func HTTPInterceptor(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		userName := req.Form.Get("username")
		token := req.Form.Get("token")
		if userName == "null" || token == "null" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ifTokenIsValid(userName, token) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		handlerFunc(w, req)
	}
}
