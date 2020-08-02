package main

import (
	"FILESTORE-SERVER/handler"
	"fmt"
	"net/http"
)

/*
	文件的校a验
	校验算法类型：CRC(32/64) MD5 SHA1
	从各方面来评判哪种算法合适： 校验值长度、校验值类型、安全级别、计算效率、应用场景
*/

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadFileSucHandler)
	http.HandleFunc("/file/metainfo", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.QueryFileMetasHandler)
	http.HandleFunc("/file/download", handler.FileDownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileMetaDeleteHandler)
	http.HandleFunc("/file/downloadurl", handler.DownloadURLHandler)
	// 秒传接口
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	// 初始化分块上传接口
	http.HandleFunc("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	// 分块上传接口
	http.HandleFunc("/file/mpupload/uppart", handler.UploadChunkFileHandler)
	// 完成分块上传接口
	http.HandleFunc("/file/mpupload/complete", handler.CompleteUploadHandler)
	// 取消文件上传接口
	http.HandleFunc("/file/mpupload/cancel", handler.CancelUploadHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SigninHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, err: %v\n", err)
	}
}