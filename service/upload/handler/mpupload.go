package handler

import (
	cacheLayer "FILESTORE-SERVER/cache/redis"
	"FILESTORE-SERVER/common"
	"FILESTORE-SERVER/mq"
	dbCli "FILESTORE-SERVER/service/dbproxy/client"
	"FILESTORE-SERVER/service/upload/config"
	"FILESTORE-SERVER/service/upload/proto"
	"FILESTORE-SERVER/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	tmpStoreDir = "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/"
	chunkSize = 1024 * 1024		// 1MB
	hashUpIdPrefixKey = "HASH_UPID_"
)

func init() {
	if err := os.MkdirAll(tmpStoreDir, 0744); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type MultipartUploadInfo struct {
	UploadID string
	FileHash string
	FileSize int
	ChunkSize int
	ChunkCount int
	ChunkExists []int	// 用于记录已上传了哪些分块
}

func (u *Upload) InitialMultipartUpload(ctx context.Context, req *proto.ReqInitialMultipartUpload, resp *proto.RespInitialMultipartUpload) error {
	//1. 解析请求参数，包括文件哈希值、文件大小、用户名
	userName := req.Username
	fileHash := req.Filehash
	fileSize := req.Filesize
	//2. 尝试获取一个redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	conn := connectionPool.Get()
	defer conn.Close()
	if err != nil {
		resp.Code = -2
		resp.Msg = "Try to multipart upload file failed, please check log to get more details!"
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return err
	}

	//2.1 通过文件哈希码与upload id对应的key，到缓存中查找当前请求携带的文件哈希码是否有对应的upload id
	//2.1.1 如果有，则说明本次上传是断点续传
	//2.1.2 如果没有，则说明本次上传是全新的分块上传
	var chunkExists []int
	tmpResult := float64(fileSize) / chunkSize  // golang中也存在隐式转换
	uploadIdInCache, _ := redis.String(conn.Do("GET", hashUpIdPrefixKey+fileHash))
	if uploadIdInCache == "" {
		chunkExists = []int{}
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
	multipartUploadInfo := MultipartUploadInfo{
		FileHash: fileHash,
		FileSize: int(fileSize),
		ChunkSize: chunkSize,
		ChunkCount: int(math.Ceil(tmpResult)),
		ChunkExists: chunkExists,
	}
	if uploadIdInCache != "" {
		multipartUploadInfo.UploadID = uploadIdInCache
	} else {
		multipartUploadInfo.UploadID = userName + fmt.Sprintf("%x", time.Now().UnixNano())
	}
	multipartUploadInfoJsonStr, err := json.Marshal(multipartUploadInfo)
	if err != nil {
		resp.Code = -2
		resp.Msg = "Try to multipart upload file failed, please check log to get more details!"
		log.Printf("%v\n", err)
		return err
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
	resp.Code = 0
	resp.Msg = "Succeed"
	resp.Initialresult = multipartUploadInfoJsonStr
	return nil
}

func (u *Upload) UploadChunkFile(ctx context.Context, req *proto.ReqUploadChunkFile, resp *proto.RespUploadChunkFile) error {
	//1. 解析请求获得参数，包括uploadid，chunk index，chunk hash
	uploadId := req.Uploadid
	chkIndex := strconv.Itoa(int(req.Chkidx))
	//2. 获得redis连接
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		resp.Code = -2
		resp.Msg = "Upload chunk file error, please check log to get more details"
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return err
	}
	conn := connectionPool.Get()
	defer conn.Close()
	//3. 本地创建相应的(分块)文件句柄，用来持久化此次客户端上传的分块文件内容
	fPath := tmpStoreDir + uploadId + "/" + chkIndex
	err = os.MkdirAll(path.Dir(fPath), 0744)
	if err != nil {
		resp.Code = -2
		resp.Msg = "Upload chunk file error, please check log to get more details"
		log.Printf("Create tmp chunk file store location failed: %v\n", err)
		return err
	}
	file, err := os.Create(fPath)
	if err != nil {
		resp.Code = -2
		resp.Msg = "Upload chunk file error, please check log to get more details"
		log.Printf("Create tmp chunk file store location failed: %v\n", err)
		return err
	}
	defer file.Close()
	file.Write(req.Chkcontent)
	//4. 更新缓存中此次分块文件所对应的分块上传信息
	conn.Do("HSET", "MP_" + uploadId, "chkidx_" + chkIndex, 1)
	//5. 返回操作成功的信息
	resp.Code = 0
	resp.Msg = "Succeed"
	return nil
}

func (u *Upload) CompleteMultipartUpload(ctx context.Context, req *proto.ReqCompleteMultipartUpload, resp *proto.RespCompleteMultipartUpload) error {
	//1. 获取请求参数，包括upload id，username，filehash，filesize，filename
	uploadId := req.Uploadid
	userName := req.Username
	fileHash := req.Filehash
	fileSize := int(req.Filesize)
	fileName := req.Filename
	//2. 获取连接池的连接，取出upload id对应的所有文件上传信息
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return err
	}
	conn := connectionPool.Get()
	defer conn.Close()
	data, err :=  redis.Values(conn.Do("HGETALL", "MP_"+uploadId))
	if err != nil {
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Printf("error: %v\n", err)
		return err
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
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Printf("The chunk file count are mismatch in two side!\n")
		panic("The chunk file count are mismatch in two side!")
	}
	//3. 合并分块文件
	srcPath := tmpStoreDir + uploadId + "/"
	destPath := tmpStoreDir + fileName + ".jpg"
	cmd := fmt.Sprintf("cd %s && cat $(ls | sort -n) > %s", srcPath, destPath)
	executeResult, err := utils.ExecuteShell(cmd)
	if err != nil {
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Println(err)
		return err
	}
	log.Printf("Merge file chunk succeed: %v\n", executeResult)
	//转移至OSS，如若有失败的情况，则进行重试
	transferData := mq.TransferData{
		FileHash: fileHash,
		CurLocation: destPath,
		DestLocation: config.OssPrefixPath + fileName,
		DestStoreType: common.StoreOSS,
	}
	bytes, _ := json.Marshal(transferData)
	processSuc := mq.Publish(mq.Exchange, mq.RoutingKey, bytes)
	if !processSuc {
		log.Println("MQ publish transferdata failed, will retry in future......")
		//todo 如果publish失败，则进行重试 2020.09.03
	}
	//4. 对唯一文件表与用户文件表新插入一条记录
	_, err = dbCli.OnFileUploadFinished(fileHash, fileName, int64(fileSize), destPath)
	if err != nil {
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Println(err)
		return err
	}
	_, err = dbCli.OnUserFileUploadFinish(userName, fileName, fileHash, int64(fileSize))
	if err != nil {
		resp.Code = -2
		resp.Msg = "Post operation failed after complete upload file, please check log!"
		log.Println(err)
		return err
	}
	//4. 返回操作成功信息
	resp.Code = 0
	resp.Msg = "Succeed"
	return nil
}

func (u *Upload) CancelUpload(ctx context.Context, req *proto.ReqCancelUpload, resp *proto.RespCancelUpload) error {
	//1. 解析请求参数，获得文件的哈希值
	fileHash := req.Filehash
	//2. 根据文件哈希值查询缓存中是否存在相对应的upload id
	connectionPool, err := cacheLayer.GetRedisConnectionPool()
	if err != nil {
		resp.Code = -2
		resp.Msg = "Cancel upload failed, please check log to get more details"
		log.Printf("Get redis connection failed: %v, please check!\n", err)
		return err
	}
	conn := connectionPool.Get()
	defer conn.Close()
	uploadId, err := redis.String(conn.Do("GET", hashUpIdPrefixKey+fileHash))
	//2.1 如果没有，则返回相关提示信息
	if err != nil || uploadId == "" {
		resp.Code = -2
		resp.Msg = "Cancel upload failed, please check log to get more details"
		log.Printf("Cancel upload file failed: %v, please check log.\n", err)
		return err
	}
	//2.2 如果有，则删除缓存中对应的信息，以及已上传的文件(删除指定目录)
	_, delHashUpIdErr := conn.Do("DEL", hashUpIdPrefixKey+fileHash)
	_, delInitialUpInfoErr := conn.Do("DEL", "MP_"+uploadId)
	if delHashUpIdErr != nil || delInitialUpInfoErr != nil {
		resp.Code = -2
		resp.Msg = "Cancel upload failed, please check log to get more details"
		log.Fatal(delHashUpIdErr)
		log.Fatal(delInitialUpInfoErr)
		return delHashUpIdErr
	}
	execResult := utils.RemovePathByShell(tmpStoreDir + uploadId)
	if !execResult {
		resp.Code = -2
		resp.Msg = "Cancel upload failed, please check log to get more details"
		panic("Can not remove tmp file, please check")
	}
	resp.Code = 0
	resp.Msg = "Succeed"
	return nil
}