package handler

import (
	cacheLayer "FILESTORE-SERVER/cache/redis"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	tmpStoreDir = "/Users/behe/Desktop/work_station/FILESTORE-SERVER/tmp/"
	chunkSize = 1024 * 1024
	hashUpIdPrefixKey = "HASH_UPID_"
)

type MultipartUploadInfo struct {
	UploadID string
	FileHash string
	FileSize int
	ChunkSize int
	ChunkCount int
	ChunkExists []int	// 用于记录已上传了哪些分块
}

func init() {
	// 创建存放分块文件的临时目录
	if err := os.MkdirAll("/Users/behe/Desktop/work_station/FILESTORE-SERVER/tmp/", 0744); err != nil {
		log.Fatal(err)
	}
}

// 初始化分块上传接口
func InitialMultipartUploadHandler(c *gin.Context) {
	//1. 解析请求参数，包括文件哈希值、文件大小、用户名
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "Invalid request parameter, please check log to get more details!",
		})
		log.Printf("Invalid request parameter: filesize, please check!")
		return
	}
	//2. 尝试获取一个redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	conn := connectionPool.Get()
	defer conn.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Try to multipart upload file failed, please check log to get more details!",
		})
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return
	}

	//2.1 通过文件哈希码与upload id对应的key，到缓存中查找当前请求携带的文件哈希码是否有对应的upload id
	//2.1.1 如果有，则说明本次上传是断点续传
	//2.1.2 如果没有，则说明本次上传是全新的分块上传
	var chunkExists []int

	uploadIdInCache, _ := redis.String(conn.Do("GET", hashUpIdPrefixKey+fileHash))
	if uploadIdInCache == "" {

	} else {
		data, _ := redis.Values(conn.Do("HGETALL", "MP_"+uploadIdInCache))
		for i := 0; i < len(data); i += 2 {
			k := string(data[i].([]byte))
			v := string(data[i + 1].([]byte))
			if strings.HasPrefix(k, "chkidx_") && v == "1" {
				chkIdx, _ := strconv.Atoi(k[7:])
				chunkExists = append(chunkExists, chkIdx)
			}
		}
	}

	//3. 构造初始化信息
	tmpResult := float64(fileSize) / chunkSize  // golang中也存在隐式转换
	multipartUploadInfo := MultipartUploadInfo{
		UploadID: userName + fmt.Sprintf("%x", time.Now().UnixNano()),
		FileHash: fileHash,
		FileSize: fileSize,
		ChunkSize: chunkSize,
		ChunkCount: int(math.Ceil(tmpResult)),
		ChunkExists: chunkExists,
	}
	multipartUploadInfoJsonStr, err := json.Marshal(multipartUploadInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Try to multipart upload file failed, please check log to get more details!",
		})
		log.Printf("%v\n", err)
		return
	}
	if len(chunkExists) == 0 {
		//4. 将初始化信息存储进redis，同时对每种信息的key值进行过期处理
		hkey := "MP_" + multipartUploadInfo.UploadID
		conn.Do("HSET", hkey, "filehash", fileHash)
		conn.Do("HSET", hkey, "filesize", fileSize)
		conn.Do("HSET", hkey, "chunkcount", multipartUploadInfo.ChunkCount)
		conn.Do("EX", hkey, 43200)
		conn.Do("SET", hashUpIdPrefixKey + fileHash, multipartUploadInfo.UploadID, "EX", 43200)
	}
	//5. 将此次操作成功的结果返回给客户端
	c.Data(http.StatusOK, "application/json", multipartUploadInfoJsonStr)
}

// 分块上传接口
func UploadChunkFileHandler(c *gin.Context) {
	//1. 解析请求获得参数，包括uploadid，chunk index，chunk hash
	uploadId := c.Request.FormValue("uploadid")
	chkIndex := c.Request.FormValue("index")
	//2. 获得redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload chunk file error, please check log to get more details",
		})
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return
	}
	conn := connectionPool.Get()
	defer conn.Close()
	//3. 本地创建相应的(分块)文件句柄，用来持久化此次客户端上传的分块文件内容
	fPath := tmpStoreDir + uploadId + "/" + chkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	file, err := os.Create(fPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Upload chunk file error, please check log to get more details",
		})
		log.Printf("Create tmp chunk file store location failed: %v\n", err)
		return
	}
	defer file.Close()
	buf := make([]byte, chunkSize) // 1MB
	for {

		n, err := c.Request.Body.Read(buf)  // 读到文件最后结束时会遇到EOF，于是会抛出err
		file.Write(buf[:n])
		if err != nil {
			fmt.Printf("Read content from request body failed: %v\n", err)
			break
		}
	}
	//4. 更新缓存中此次分块文件所对应的分块上传信息
	conn.Do("HSET", "MP_" + uploadId, "chkidx_" + chkIndex, 1)
	//5. 返回操作成功的信息
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Upload chunk file succeed!",
	})
}

// 完成分块上传接口
func CompleteUploadHandler(c *gin.Context) {
	//1. 获取请求参数，包括upload id，username，filehash，filesize，filename
	uploadId := c.Request.FormValue("uploadid")
	userName := c.Request.FormValue("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize, _ := strconv.Atoi(c.Request.FormValue("filesize"))
	fileName := c.Request.FormValue("filename")
	//2. 获取连接池的连接，取出upload id对应的所有文件上传信息
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Post operation failed after complete upload file, please check log!",
		})
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return
	}
	conn := connectionPool.Get()
	defer conn.Close()
	data, err :=  redis.Values(conn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Post operation failed after complete upload file, please check log!",
		})
		log.Printf("error: %v\n", err)
		return
	}
	exceptCount := 0
	actualCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i + 1].([]byte))
		if k == "chunkcount" {
			exceptCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1"{
			actualCount++
		}
	}
	if exceptCount != actualCount {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Post operation failed after complete upload file, please check log!",
		})
		log.Printf("The chunk file count are mismatch in two side!\n")
		return
	}
	//3. 验证文件上传是否完整，如果确实上传完整，则对唯一文件表与用户文件表新插入一条记录
	db.OnFileUploadFinished(fileHash, fileName, int64(fileSize), "")
	db.OnUserFileUploadFinish(userName, fileName, fileHash, int64(fileSize))
	//4. 返回操作成功信息
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Post operation succeed!",
	})
}

// 取消文件上传接口
func CancelUploadHandler(c *gin.Context) {
	//1. 解析请求参数，获得文件的哈希值
	fileHash := c.Request.FormValue("filehash")
	//2. 根据文件哈希值查询缓存中是否存在相对应的upload id
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Cancel upload failed, please check log to get more details",
		})
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return
	}
	conn := connectionPool.Get()
	defer conn.Close()
	uploadId, err := redis.String(conn.Do("GET", hashUpIdPrefixKey+fileHash))
	//2.1 如果没有，则返回相关提示信息
	if err != nil || uploadId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Cancel upload failed, please check log to get more details",
		})
		log.Printf("Cancel upload file failed: %v, please check log.\n", err)
		return
	}
	//2.2 如果有，则删除缓存中对应的信息，以及已上传的文件(删除指定目录)
	_, delHashUpIdErr := conn.Do("DEL", hashUpIdPrefixKey+fileHash)
	_, delInitialUpInfo := conn.Do("DEL", "MP_"+uploadId)
	if delHashUpIdErr != nil || delInitialUpInfo != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Cancel upload failed, please check log to get more details",
		})
		log.Fatal(delHashUpIdErr)
		log.Fatal(delInitialUpInfo)
		return
	}
	execResult := utils.RemovePathByShell(tmpStoreDir + uploadId)
	if !execResult {
		c.JSON(http.StatusOK, gin.H{
			"code": -2,
			"msg": "Cancel upload failed, please check log to get more details",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": "Cancel upload succeed!",
	})
}
