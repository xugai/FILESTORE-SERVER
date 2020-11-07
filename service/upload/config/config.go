package config

import "FILESTORE-SERVER/common"

var UploadEntry = "localhost:30082"		// 配置上传入口的地址
var UploadServiceHost = ":28080"	// 上传服务监听的地址

var CurrentStoreType = common.StoreOSS
const OssPrefixPath = "oss/image/"


