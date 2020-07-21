package test

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	userName = "behe"
	token = ""
	fileHash = "fa5e59e42718682b51d28528b3eb85b815af1e0b"
	fileSize = "4910254"
	fileName = "06753103-B569-4AE8-B282-8FFD94C98730.jpeg"
)

func TestMpUpload(t *testing.T) {

	//1. 测试初始化分块上传接口
	resp, err := http.PostForm("http://localhost:8080/file/mpupload/init", url.Values{
		"username": {userName},
		"filehash": {fileHash},
		"filesize": {fileSize},
		"token":    {token},
	})
	if err != nil {
		t.Fatal(err)
		//log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// 1.1 获得服务端返回回来的upload id与chunk size
	uploadId := jsoniter.Get(body, "Data").Get("UploadID").ToString()
	chunkSize := jsoniter.Get(body, "Data").Get("ChunkSize").ToInt()
	fmt.Printf("Get uploadId: %v, chunkSize: %v\n", uploadId, chunkSize)

	//2. 测试文件分块上传接口
	err = UpMultipart(uploadId, chunkSize)
	if err != nil {
		t.Fatal(err)
	}
	// 3. 测试文件分块上传完成接口
	response, err := http.PostForm("http://localhost:8080/file/mpupload/complete", url.Values{
		"uploadid": {uploadId},
		"username": {userName},
		"filehash": {fileHash},
		"filename": {fileName},
		"filesize": {fileSize},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, _ = ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

}
