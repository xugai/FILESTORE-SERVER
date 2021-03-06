package mpupload

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
	fileHash = "1eeb7d8383896ec89be67b3573af94952a1d726f"
	fileSize = "196707"
	fileName = "BAEE.jpg"
)

func convertInterfaceArrToIntArr(interfaceArr []interface{}) []int {
	var intArr []int
	for i := 0; i < len(interfaceArr); i++ {
		intArr = append(intArr, int(interfaceArr[i].(float64)))
	}
	return intArr
}

func TestMpUpload(t *testing.T) {

	//1. 测试初始化分块上传接口
	resp, err := http.PostForm("http://localhost:28080/file/mpupload/init", url.Values{
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
	uploadId := jsoniter.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsoniter.Get(body, "data").Get("ChunkSize").ToInt()
	chunkExists := convertInterfaceArrToIntArr(
							jsoniter.Get(body, "data").
							Get("ChunkExists").
							GetInterface().([]interface{}))
	fmt.Printf("Get uploadId: %v, chunkSize: %v\n", uploadId, chunkSize)

	//2. 测试文件分块上传接口
	//2.1 如果已上传的分块数量为0，则说明之前从未上传过，因此这次是分块上传
	//2.2 否则这次就是断点续传
	if len(chunkExists) == 0 {
		err = UpMultipart(uploadId, chunkSize)
	} else if len(chunkExists) > 0 {
		err = ResumeBreakpoint(uploadId, chunkSize, chunkExists)
	}
	if err != nil {
		t.Fatal(err)
	}
	// 3. 测试文件分块上传完成接口
	response, err := http.PostForm("http://localhost:28080/file/mpupload/complete", url.Values{
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
