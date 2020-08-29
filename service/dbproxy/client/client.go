package client

import (
	"FILESTORE-SERVER/service/dbproxy/mapper"
	"FILESTORE-SERVER/service/dbproxy/proto"
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/mitchellh/mapstructure"
	"log"
)

type FileMeta struct {
	UploadAt string
	FileName string
	FileSha1 string
	FileSize int64
	Location string
}

var dbCli proto.DBProxyService

func init() {
	newRegistry := consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	service := micro.NewService(
		micro.Registry(newRegistry),
		)
	service.Init()
	dbCli = proto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

func TableFileToFileMeta(file mapper.TableFile) FileMeta {
	return FileMeta{
		FileName: file.FileName.String,
		FileSha1: file.FileHash.String,
		FileSize: file.FileSize.Int64,
		Location: file.FileAddr.String,
	}
}

// 向rpc server请求执行action
func execAction(name string, paramJson []byte) (*proto.RespExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &proto.ReqExec{
		Action: []*proto.SingleAction{
			{
				Name: name,
				Params: paramJson,
			},
		},
	})
}

// 解析rpc server返回的data
func parseBody(resp *proto.RespExec) *mapper.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	results := []mapper.ExecResult{}
	err := json.Unmarshal(resp.Data, &results)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(results) > 0 {
		return &results[0]
	}
	return nil
}

func ToFileMeta(src interface{}) FileMeta {
	fileMeta := FileMeta{}
	mapstructure.Decode(src, &fileMeta)
	return fileMeta
}

func ToTableUser(src interface{}) mapper.User {
	user := mapper.User{}
	mapstructure.Decode(src, &user)
	return user
}

func ToTableFile(src interface{}) mapper.TableFile {
	file := mapper.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func ToTableUserFile(src interface{}) mapper.UserFile {
	userFile := mapper.UserFile{}
	mapstructure.Decode(src, &userFile)
	return userFile
}

func ToTableUserFiles(src interface{}) []mapper.UserFile {
	userFiles := []mapper.UserFile{}
	mapstructure.Decode(src, &userFiles)
	return userFiles
}

func GetFileMeta(filehash string) (*mapper.ExecResult, error){
	bytes, _ := json.Marshal([]interface{}{filehash})
	respExec, err := execAction("/file/GetFileMeta", bytes)
	return parseBody(respExec), err
}

func OnFileUploadFinished(fileHash string, fileName string, fileSize int64, fileAddr string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{fileHash, fileName, fileSize, fileAddr})
	respExec, err := execAction("/file/OnFileUploadFinished", bytes)
	return parseBody(respExec), err
}

func UpdateFileStoreLocation(fileHash, fileAddr string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{fileHash, fileAddr})
	respExec, err := execAction("/file/UpdateFileStoreLocation", bytes)
	return parseBody(respExec), err
}

func UserSignup(userName string, passWord string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, passWord})
	respExec, err := execAction("/user/UserSignup", bytes)
	return parseBody(respExec), err
}

func UserSignin(userName string, passWord string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, passWord})
	respExec, err := execAction("/user/UserSignin", bytes)
	return parseBody(respExec), err
}

func FlushUserToken(userName string, token string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, token})
	respExec, err := execAction("/user/FlushUserToken", bytes)
	return parseBody(respExec), err
}

func GetUserInfo(userName string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName})
	respExec, err := execAction("/user/GetUserInfo", bytes)
	return parseBody(respExec), err
}

func IfTokenIsValid(username string, token string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{username, token})
	respExec, err := execAction("/user/IfTokenIsValid", bytes)
	return parseBody(respExec), err
}

func OnUserFileUploadFinish(userName string, fileName string, fileHash string, fileSize int64) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, fileName, fileHash, fileSize})
	respExec, err := execAction("/ufile/OnUserFileUploadFinish", bytes)
	return parseBody(respExec), err
}

func GetUserFileMetas(userName string, limit int) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, limit})
	respExec, err := execAction("/ufile/GetUserFileMetas", bytes)
	return parseBody(respExec), err
}

func GetUserFileMeta(userName string, fileHash string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, fileHash})
	respExec, err := execAction("/ufile/GetUserFileMeta", bytes)
	return parseBody(respExec), err
}

func UpdateUserFileMeta(userName string, fileHash string, newFileName string) (*mapper.ExecResult, error) {
	bytes, _ := json.Marshal([]interface{}{userName, fileHash, newFileName})
	respExec, err := execAction("/ufile/UpdateUserFileMeta", bytes)
	return parseBody(respExec), err
}




