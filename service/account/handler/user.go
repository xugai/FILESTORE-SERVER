package handler

import (
	"FILESTORE-SERVER/config"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/service/account/proto"
	"FILESTORE-SERVER/utils"
	"context"
)

type User struct {}

// 处理用户注册请求
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	userName := req.Username
	passWord := req.Password
	// validation
	if len(userName) < 3 || len(passWord) < 5 {
		resp.Code = -1
		resp.Message = "Invalid request parameter(s)!"
		return nil
	}
	// encrypt for password
	result := db.UserSignup(userName, utils.Sha1([]byte(passWord + config.Salt)))
	if !result {
		resp.Code = -2
		resp.Message = "Sign up failed!"
		return nil
	}
	resp.Code = 0
	resp.Message = "Sign up succeed!"
	return nil
}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {
	userName := req.Username
	passWord := req.Password

	result := db.UserSignin(userName, utils.Sha1([]byte(passWord + config.Salt)))
	if !result {
		resp.Code = -2
		resp.Message = "Login failed, your username or password maybe is incorrect!"
		return nil
	}
	token := utils.GenToken(userName)
	result = db.FlushUserToken(userName, token)
	if !result {
		resp.Code = -2
		resp.Message = "Flush user token failed, please check log"
		return nil
	}
	resp.Code = 0
	resp.Token = token
	resp.Message = "Login succeed"
	return nil
}

func (u *User) 	UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {
	userName := req.Username
	userInfo, err := db.GetUserInfo(userName)
	if err != nil {
		resp.Code = -2
		resp.Message = "Get user info failed, please check!"
		return err
	}
	resp.Code = 0
	resp.Message = "OK"
	resp.Username = userInfo.UserName
	resp.Email = userInfo.Email
	resp.Phone = userInfo.Phone
	resp.LastActiveAt = userInfo.LastActive
	resp.Signup = userInfo.SignupAt
	resp.Status = int32(userInfo.Status)
	return nil
}
