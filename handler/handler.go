package handler

import (
	"FILESTORE-SERVER/common"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/mq"
	"FILESTORE-SERVER/store/oss"
	"FILESTORE-SERVER/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var currentStoreType = common.StoreOSS
const ossPrefixPath = "oss/image/"

func UploadHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/index.html")
}

// for gin web fwk
func UploadPostHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	userName := c.Request.FormValue("username")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Failed to get uploaded file data, please check log to get more details!",
		})
		log.Printf("Failed to get uploaded file data, err: %v\n", err)
		return
	}
	defer file.Close()
	fileMeta := meta.FileMeta{
		Location: "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/" + header.Filename,
		FileName: header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload file failed, please check log to get more details!",
		})
		log.Printf("Failed to create file in /tmp/" + header.Filename +", err: %v\n", err)
		return
	}
	defer newFile.Close()
	fileMeta.FileSize, err = io.Copy(newFile, file)
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = utils.FileSha1(newFile)
	log.Printf("%v Upload file with hash: %v", time.Now().Format("2006-01-02 15:04:05"), fileMeta.FileSha1)

	// 将文件以同步/异步的方式转移到公有云OSS上
	newFile.Seek(0, 0)
	if common.StoreOSS == currentStoreType {
		ossObjectKey := ossPrefixPath + fileMeta.FileName
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
				c.JSON(http.StatusOK, gin.H{
					"code": -2,
					"msg": "Upload file failed, please check log to get more details!",
				})
				log.Printf("%v\n", err)
				return
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
				c.JSON(http.StatusOK, gin.H{
					"code": -2,
					"msg": "Upload file failed, please check log to get more details!",
				})
				log.Printf("Put object to Ali OSS failed: %v\n", err)
				return
			}
		}
	}

	result := meta.UpdateFileMetaDB(fileMeta)
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload file failed, please check log to get more details!",
		})
		log.Printf("Update or insert file meta in DB error, please check.\n")
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload file failed, please check log to get more details!",
		})
		log.Printf("Failed to save data to new file, err: %v\n",  err)
		return
	}
	result = db.OnUserFileUploadFinish(userName, fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize)
	if result {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg": "Upload file succeed!",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload file failed, please check log to get more details!",
		})
	}
}

func UploadFileSucHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Upload finish",
	})
}

func GetFileMetaHandler(c *gin.Context){
	filehash := c.Request.FormValue("filehash")
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err2 := meta.GetFileMetaDB(filehash)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Get file meta error, please check log to get more details!",
		})
		log.Printf("Get file meta error occured: %v\n", err2)
		return
	}

	fileMetaInfo, err := json.Marshal(*fileMeta)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Get file meta error, please check log to get more details!",
		})
		log.Printf("%v\n", err)
		return
	}
	c.Data(http.StatusOK, "application/json", fileMetaInfo)
}

func QueryFileMetasHandler(c *gin.Context) {
	limitCnt, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg": "Request parameters are invalid!",
		})
		return
	}
	userName := c.Request.FormValue("username")
	userFiles, _ := db.GetUserFileMetas(userName, limitCnt)
	fileMetaArrayJsonStr, err := json.Marshal(userFiles)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Query file meta info failed, please check log to get more details!",
		})
		log.Printf("%v\n", err)
		return
	}
	c.Data(http.StatusOK, "application/json", fileMetaArrayJsonStr)
}

func FileDownloadHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	userName := c.Request.FormValue("username")
	fm := meta.GetFileMeta(fileHash)
	userFile, err := db.GetUserFileMeta(userName, fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Download file error, please check log to get more details!",
		})
		return
	}
	if !strings.HasPrefix(fm.Location, ossPrefixPath) {
		// 直接从本地下载返回
		c.FileAttachment(fm.Location, userFile.FileName)
	} else {
		// 从oss上下载
		downloadURL := oss.Download(fm.Location)
		request, err := http.NewRequest(http.MethodPost, downloadURL, nil)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -2,
				"msg": "Download file error, please check log to get more details!",
			})
			log.Printf("%v\n", err)
			return
		}
		defer request.Body.Close()
		readAll, _ := ioutil.ReadAll(request.Body)
		c.Data(http.StatusOK, "application/octect-stream", readAll)
	}
}

func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	if opType != "0" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg": "Request parameters invalid!",
		})
		return
	}
	userName := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")

	if !db.UpdateUserFileMeta(userName, filehash, newFileName) {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "File meta info update failed, please check log to get more details!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "File meta info update succeed",
	})
}

func FileMetaDeleteHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")
	err := os.Remove(meta.GetFileMeta(filehash).Location)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Delete file error, please check log to get more details!",
		})
		log.Printf("Delete file err: %v\n", err)
		return
	}
	meta.RemoveFileMeta(filehash)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Delete file meta succeed",
	})
}

// 尝试秒传处理用户上传的文件
func TryFastUploadHandler(c *gin.Context) {
	// 获取页面传回来的用户名、当前上传文件的SHA1值
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	fileSize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	// 根据SHA1值在tbl_file中查找是否已经存在，如果存在，则往tbl_user_file中插入新的记录，否则秒传失败
	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Fast upload file error, please check log to get more details!",
		})
		log.Printf("Get file meta error: %v\n", err)
		return
	}
	if fileMeta == nil {
		//todo 如果云端中确实没有已上传的文件，是不是这里应该主动跳转为普通上传？
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Fast upload failed, please use regular upload",
		})
		return
	}
	result := db.OnUserFileUploadFinish(userName, fileName, fileMeta.FileHash.String, int64(fileSize))
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Fast upload file error, please check log to get more details!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Fast upload file succeed",
	})
}