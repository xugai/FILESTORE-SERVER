package main

import (
	"FILESTORE-SERVER/config"
	"FILESTORE-SERVER/route"
)

/*
	文件的校a验
	校验算法类型：CRC(32/64) MD5 SHA1
	从各方面来评判哪种算法合适： 校验值长度、校验值类型、安全级别、计算效率、应用场景
*/

func main() {

	//http.Handle("/static/",
	//	http.StripPrefix("/static/",
	//		http.FileServer(http.Dir("./static"))))
	router := route.Router()
	router.Run(config.UploadServiceHost)
}

