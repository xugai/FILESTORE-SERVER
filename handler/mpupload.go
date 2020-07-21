package handler

import (
	cacheLayer "FILESTORE-SERVER/cache/redis"
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/utils"
	"fmt"
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
)

type MultipartUploadInfo struct {
	UploadID string
	FileHash string
	FileSize int
	ChunkSize int
	ChunkCount int
}

func init() {
	if err := os.MkdirAll("/Users/behe/Desktop/work_station/FILESTORE-SERVER/tmp/", 0744); err != nil {
		log.Fatal(err)
	}
}

func InitialMultipartUploadHandler(w http.ResponseWriter, req *http.Request) {
	//1. 解析请求参数，包括文件哈希值、文件大小、用户名
	req.ParseForm()
	userName := req.Form.Get("username")
	fileHash := req.Form.Get("filehash")
	fileSize, err := strconv.Atoi(req.Form.Get("filesize"))
	if err != nil {
		fmt.Printf("Invalid request parameter: filesize, please check!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.NewSimpleServerResponse(400, "请求参数非法，请检查").GetInByteStream())
		return
	}
	//2. 尝试获取一个redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		fmt.Printf("Get redis connection failed: %v, please check!\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewSimpleServerResponse(500, "服务内部发生异常错误，请检查活动日志!").GetInByteStream())
		return
	}

	//3. 构造初始化信息
	tmpResult := float64(fileSize) / chunkSize  // golang中也存在隐式转换
	multipartUploadInfo := MultipartUploadInfo{
		UploadID: userName + fmt.Sprintf("%x", time.Now().UnixNano()),
		FileHash: fileHash,
		FileSize: fileSize,
		ChunkSize: chunkSize,
		ChunkCount: int(math.Ceil(tmpResult)),
	}
	//4. 将初始化信息存储进redis
	conn := connectionPool.Get()
	defer conn.Close()
	conn.Do("HSET", "MP_" + multipartUploadInfo.UploadID, "filehash", fileHash)
	conn.Do("HSET", "MP_" + multipartUploadInfo.UploadID, "filesize", fileSize)
	conn.Do("HSET", "MP_" + multipartUploadInfo.UploadID, "chunkcount", multipartUploadInfo.ChunkCount)

	//5. 将此次操作成功的结果返回给客户端
	w.Write(utils.NewServerResponse(200, "初始化分块上传信息成功!", multipartUploadInfo).GetInByteStream())
}

func UploadChunkFileHandler(w http.ResponseWriter, req *http.Request) {
	//1. 解析请求获得参数，包括uploadid，chunk index，chunk hash
	req.ParseForm()
	uploadId := req.Form.Get("uploadid")
	chkIndex := req.Form.Get("index")
	//2. 获得redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		fmt.Printf("Get redis connection failed: %v, please check!\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewSimpleServerResponse(500, "服务内部发生异常错误，请检查活动日志!").GetInByteStream())
		return
	}
	conn := connectionPool.Get()
	defer conn.Close()
	//3. 本地创建相应的(分块)文件句柄，用来持久化此次客户端上传的分块文件内容
	fPath := tmpStoreDir + uploadId + "/" + chkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	file, err := os.Create(fPath)
	if err != nil {
		fmt.Printf("Create tmp chunk file store location failed: %v\n", err)
		w.Write(utils.NewSimpleServerResponse(500, "创建临时分块文件存储文件失败!").GetInByteStream())
		return
	}
	defer file.Close()
	buf := make([]byte, 1024 * 1024) // 1MB
	for {
		n, err := req.Body.Read(buf)  // 读到文件最后结束时会遇到EOF，于是会抛出err
		file.Write(buf[:n])
		if err != nil {
			fmt.Printf("Read content from request body failed: %v\n", err)
			break
		}
	}
	//4. 更新缓存中此次分块文件所对应的分块上传信息
	conn.Do("HSET", "MP_" + uploadId, "chkidx_" + chkIndex, 1)
	//5. 返回操作成功的信息
	w.Write(utils.NewSimpleServerResponse(200, "分块文件上传成功!").GetInByteStream())
}

func CompleteUploadHandler(w http.ResponseWriter, req *http.Request) {
	//1. 获取请求参数，包括upload id，username，filehash，filesize，filename
	req.ParseForm()
	uploadId := req.Form.Get("uploadid")
	userName := req.Form.Get("username")
	fileHash := req.Form.Get("filehash")
	fileSize, _ := strconv.Atoi(req.Form.Get("filesize"))
	fileName := req.Form.Get("filename")
	//2. 获取连接池的连接，取出upload id对应的所有文件上传信息
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		fmt.Printf("Get redis connection failed: %v, please check!\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewSimpleServerResponse(500, "服务内部发生异常错误，请检查活动日志!").GetInByteStream())
		return
	}
	conn := connectionPool.Get()
	defer conn.Close()
	data, err :=  redis.Values(conn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		fmt.Printf("error: %v\n", err)
		w.Write(utils.NewSimpleServerResponse(500, "服务器内部发生错误，请检查日志!").GetInByteStream())
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
		fmt.Printf("The chunk file count are mismatch in two side!")
		w.Write(utils.NewSimpleServerResponse(400, "请求参数不正确，请重试!").GetInByteStream())
		return
	}
	//3. 验证文件上传是否完整，如果确实上传完整，则对唯一文件表与用户文件表新插入一条记录
	db.OnFileUploadFinished(fileHash, fileName, int64(fileSize), "")
	db.OnUserFileUploadFinish(userName, fileName, fileHash, int64(fileSize))
	//4. 返回操作成功信息
	w.Write(utils.NewSimpleServerResponse(200, "上传合并分块文件成功!").GetInByteStream())
}
