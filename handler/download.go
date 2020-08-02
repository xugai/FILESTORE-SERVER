package handler

import (
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/store/oss"
	"FILESTORE-SERVER/utils"
	"net/http"
)

const (
	objectKeyPrefix = "oss/image/"
)

func DownloadURLHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	//1. 获得用户传过来的filehash
	filehash := req.Form.Get("filehash")
	//2. 通过filehash去db里面查询相应的文件的path
	fileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		w.Write(utils.NewSimpleServerResponse(500, "系统内部发生错误，请检查活动日志!").GetInByteStream())
		return
	}
	//3. 然后通过filepath获得ali oss的signed download url
	signedURL := oss.Download(objectKeyPrefix + fileMeta.FileName.String)
	//4. 最后将signed url返回给用户
	w.Write([]byte(signedURL))
}
