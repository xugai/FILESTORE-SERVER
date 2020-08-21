package handler

import (
	"FILESTORE-SERVER/common"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/mq"
	"FILESTORE-SERVER/service/upload/config"
	"FILESTORE-SERVER/service/upload/proto"
	"FILESTORE-SERVER/store/oss"
	"FILESTORE-SERVER/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

type Upload struct {
}

func (u *Upload) UploadEntry(ctx context.Context, req *proto.ReqUploadEntry, resp *proto.RespUploadEntry) error {
	resp.Code = 0
	resp.Message = "OK"
	resp.Entry = config.UploadEntry
	return nil
}

func (u *Upload) UploadFile(ctx context.Context, req *proto.ReqUploadFile, resp *proto.RespUploadFile) error {
	fileMeta := meta.FileMeta{
		Location: "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/" + req.Filename,
		FileName: req.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		resp.Code = -2
		resp.Message = "Upload file failed, please check log to get more details!"
		log.Printf("Failed to create file in /tmp/" + req.Filename +", err: %v\n", err)
		return err
	}
	defer newFile.Close()
	// 如果发生了错误，但错误日志没有任何信息，可能说明原文件流没有被完全写入
	n, err := newFile.Write(req.Filecontent)
	if err != nil || n != len(req.Filecontent) {
		resp.Code = -2
		resp.Message = "Upload file failed, please check log to get more details!"
		log.Println(err)
	}
	newFile.Seek(0, 0)	// 游标重新回到文件头部
	fileMeta.FileSha1 = utils.FileSha1(newFile)
	log.Printf("%v Upload file with hash: %v", time.Now().Format("2006-01-02 15:04:05"), fileMeta.FileSha1)

	// 将文件以同步/异步的方式转移到公有云OSS上
	newFile.Seek(0, 0)
	if common.StoreOSS == config.CurrentStoreType {
		ossObjectKey := config.OssPrefixPath + fileMeta.FileName
		if mq.AsyncTransferEnable {
			// 异步转移文件
			transferData := mq.TransferData{
				FileHash:      fileMeta.FileSha1,
				CurLocation:   fileMeta.Location,
				DestLocation:  ossObjectKey,
				DestStoreType: common.StoreOSS,
			}
			bytes, err := json.Marshal(transferData)
			if err != nil {
				resp.Code = -1
				resp.Message = "Upload file failed, please check log to get more details!"
				log.Printf("%v\n", err)
				return err
			}
			processSucc := mq.Publish(mq.Exchange, mq.RoutingKey, bytes)
			if !processSucc {
				//todo 当前进行异步转移文件失败，稍后重试

			}
			//fileMeta.Location = ossObjectKey  等真正异步转移成功后，再来修改文件表中的存储位置
		} else {
			// 同步转移文件
			err := oss.Bucket().PutObject(ossObjectKey, newFile)
			if err != nil {
				resp.Code = -2
				resp.Message = "Upload file failed, please check log to get more details!"
				log.Printf("Put object to Ali OSS failed: %v\n", err)
				return err
			}
		}
	}

	result := meta.UpdateFileMetaDB(fileMeta)
	if !result {
		resp.Code = -2
		resp.Message = "Upload file failed, please check log to get more details!"
		log.Printf("Update or insert file meta in DB error, please check.\n")
		return errors.New("Update or insert file meta in DB error, please check")
	}
	if err != nil {
		resp.Code = -2
		resp.Message = "Upload file failed, please check log to get more details!"
		log.Printf("Failed to save data to new file, err: %v\n",  err)
		return err
	}
	result = db.OnUserFileUploadFinish(req.Username, fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize)
	if result {
		resp.Code = 0
		resp.Message = "Upload file succeed!"
		return nil
	} else {
		resp.Code = -2
		resp.Message = "Upload file failed"
		return errors.New("Upload file failed, please check log to get more details!")
	}

}


