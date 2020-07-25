package test

import (
	"FILESTORE-SERVER/utils"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	fileToUpload = "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/06753103-B569-4AE8-B282-8FFD94C98730.jpeg"
)

// 正常的客户端分块上传逻辑
func UpMultipart(uploadId string, chunkSize int) error {
	file, err := os.Open(fileToUpload)
	if err != nil {
		return err
	}
	defer file.Close()
	targetFile, _ := ioutil.ReadAll(file)
	file.Seek(0, 0)
	reader := bufio.NewReader(file)
	//startup go routine
	chunkCount := int(math.Ceil(float64(len(targetFile)) / float64(chunkSize)))
	c := make(chan bool)
	for i := 1; i <= chunkCount; i++ {
		if chunkCount == 1 {
			//chkHash := utils.Sha1(targetFile)
			//go func(b []byte) {
			//	response, err := http.Post(
			//		"http://localhost:8080/file/mpupload/uppart?username=admin&uploadid="+uploadId+"&index=1&chkhash="+chkHash,
			//		"multipart/form-data",
			//		bytes.NewReader(b))
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	defer response.Body.Close()
			//	respContent, _ := ioutil.ReadAll(response.Body)
			//	fmt.Println(string(respContent))
			//	c <- true
			//}(targetFile)
		} else {
			buf := make([]byte, chunkSize)
			n, err := reader.Read(buf)
			if err == io.EOF {
				log.Println("Read EOF with a file")
			} else if err != nil {
				return err
			}
			//newBuf := make([]byte, chunkSize)
			//copy(newBuf, buf)
			go func(b []byte, index int) {

				chkHash := utils.Sha1(b)
				response, err := http.Post(
					"http://localhost:8080/file/mpupload/uppart?username=admin&uploadid="+uploadId+"&index="+strconv.Itoa(index)+"&chkhash="+chkHash,
					"multipart/form-data",
					bytes.NewReader(b))
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()
				respContent, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(respContent))
				c <- true
			}(buf[:n], i)
		}
	}
	for k := 0; k < chunkCount; k++ {
		select {
		case <- c:
			fmt.Printf("Finish upload %d chunk.\n", k + 1)
		}
	}
	return nil
}

// 断点续传下客户端的上传逻辑
func ResumeBreakpoint(uploadId string, chunkSize int, chunkExists []int) error {
	file, err := os.Open(fileToUpload)
	if err != nil {
		return err
	}
	// 1. 获得文件大小，求出本次上传需要将文件分多少块上传
	readAll, err := ioutil.ReadAll(file)
	file.Seek(0, 0)
	if err != nil {
		return err
	}
	chunkCount := int(math.Ceil(float64(len(readAll)) / float64(chunkSize)))
	// 2. 根据服务端返回来的之前已上传过的文件分块，计算出此次剩下哪些分块是需要上传的
	chunkArray := make([]int, chunkCount)
	for i := 0; i < len(chunkExists); i++ {
		chunkArray[chunkExists[i] - 1] = 1
	}
	// 3. 找出那些还没有上传的分块，然后进行上传
	totalOfChunkToUpload := 0
	ch := make(chan bool)
	for j := 0; j < chunkCount; j++ {
		if chunkArray[j] == 0 {
			// 说明本次需要上传的文件的第j + 1块是还没有上传的
			totalOfChunkToUpload++
			if j == chunkCount - 1 {
				// 如果是文件的最后一块没有上传
				buf := readAll[j * chunkSize:]
				//todo upload last chunk with http request
				go func(b []byte, index int) {

					chkHash := utils.Sha1(b)
					response, err := http.Post(
						"http://localhost:8080/file/mpupload/uppart?username=admin&uploadid="+uploadId+"&index="+strconv.Itoa(index)+"&chkhash="+chkHash,
						"multipart/form-data",
						bytes.NewReader(b))
					if err != nil {
						log.Fatal(err)
					}
					defer response.Body.Close()
					respContent, _ := ioutil.ReadAll(response.Body)
					fmt.Println(string(respContent))
				}(buf, j + 1)
			}
			buf := readAll[j * chunkSize: (j + 1) * chunkSize]
			//todo upload chunk with http request
			go func(b []byte, index int) {
				chkHash := utils.Sha1(b)
				response, err := http.Post(
					"http://localhost:8080/file/mpupload/uppart?username=admin&uploadid="+uploadId+"&index="+strconv.Itoa(index)+"&chkhash="+chkHash,
					"multipart/form-data",
					bytes.NewReader(b))
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()
				respContent, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(respContent))
				ch <- true
			}(buf, j + 1)
		}
	}

	for k := 0; k < totalOfChunkToUpload; k++ {
		select {
		case <- ch:
			fmt.Printf("Finish upload %v chunk(s)\n", k + 1)
		}
	}

	return nil
}

