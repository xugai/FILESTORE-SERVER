package handler

import (
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/utils"
	"fmt"
	"io/ioutil"
	"net/http"
)

const salt = "s*&%#"

func SignupHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		http.Redirect(w, req, "/static/view/signup.html", http.StatusFound)
		return
	}
	req.ParseForm()
	userName := req.Form.Get("username")
	passWord := req.Form.Get("password")
	// validation
	if len(userName) < 3 || len(passWord) < 5 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request parameter(s)!"))
		return
	}
	// encrypt for password
	result := db.UserSignup(userName, utils.Sha1([]byte(passWord + salt)))
	if !result {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("FAILED"))
		return
	}
	w.Write([]byte("SUCCESS"))
}

func SigninHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		file, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			fmt.Printf("Read file signin.html error: %v\n", err)
		}
		w.Write(file)
		return
	}
	req.ParseForm()
	userName := req.Form.Get("username")
	passWord := req.Form.Get("password")

	result := db.UserSignin(userName, utils.Sha1([]byte(passWord + salt)))
	if !result {
		fmt.Printf("Login failed, please check log.")
		return
	}
	token := utils.GenToken(userName)
	result = db.FlushUserToken(userName, token)
	if !result {
		fmt.Printf("Flush user token failed, please check log.\n")
		return
	}
	serverResponse := utils.NewServerResponse(200, "OK", struct {
		Location string
		Username string
		Token    string
	}{
		Location: "http://" + req.Host + "/static/view/home.html",
		Username: userName,
		Token:    token,
	})
	w.Write(serverResponse.GetInByteStream())
}

func UserInfoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	req.ParseForm()
	userName := req.Form.Get("username")
	userInfo, err := db.GetUserInfo(userName)
	if err != nil {
		fmt.Println("Get user info failed, please check!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	serverResponse := utils.NewServerResponse(200, "OK", userInfo)
	w.Write(serverResponse.GetInByteStream())
}
