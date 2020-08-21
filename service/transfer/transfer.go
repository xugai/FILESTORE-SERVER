package main

import (
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/mq"
	"FILESTORE-SERVER/store/oss"
	"encoding/json"
	"log"
	"os"
)

var transOSSQueueName = "filestoreserver.trans.oss"
var transConsumerName = "transfer.oss"
// 消费者的callback方法
func TransferFileToOSS(msg []byte) bool {
	transferData := mq.TransferData{}
	err := json.Unmarshal(msg, &transferData)
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}
	file, err := os.Open(transferData.CurLocation)
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}
	err = oss.Bucket().PutObject(transferData.DestLocation, file)
	if err != nil {
		log.Printf("%v\n", err)
		return false
	}
	log.Printf("Succeed transfer file: %v to oss\n", transferData.FileHash)
	// 上传OSS成功后，别忘了修改db中该文件的存储位置，应修改为在OSS中的存储位置
	db.UpdateFileStoreLocation(transferData.FileHash, transferData.DestLocation)
	return true
}

// 启动消费者，开始消费
func main() {
	if !mq.AsyncTransferEnable {
		log.Printf("You do not open async transfer mode, please check your config!")
		return
	}
	log.Println("Open async transfer mode......")
	mq.DoConsume(transOSSQueueName,
				transConsumerName,
				TransferFileToOSS)
}
