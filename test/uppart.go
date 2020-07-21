package test

import (
	"FILESTORE-SERVER/utils"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

func UpMultipart(uploadId string, chunkSize int) error {
	file, err := os.Open("/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/06753103-B569-4AE8-B282-8FFD94C98730.jpeg")
	if err != nil {
		return err
	}
	defer file.Close()
	targetFile, _ := ioutil.ReadAll(file)

	//startup go routine
	chunkCount := int(math.Ceil(float64(len(targetFile)) / float64(chunkSize)))
	c := make(chan bool)
	for i := 1; i <= chunkCount; i++ {
		//get response from server
		if chunkCount == 1 {
			// no need to split
			chkHash := utils.Sha1(targetFile)
			go func(b []byte) {
				response, err := http.Post(
					"http://localhost:8080/file/mpupload/uppart?username=admin&uploadid="+uploadId+"&index=1&chkhash="+chkHash,
					"multipart/form-data",
					bytes.NewReader(targetFile))
				if err != nil {
					log.Fatal(err)
				}
				defer response.Body.Close()
				respContent, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(respContent))
				c <- true
			}(targetFile)
		} else {
			buf := make([]byte, chunkSize)
			n := 0
			if i == chunkCount {
				copy(buf, targetFile[(i - 1) * chunkSize: len(targetFile)])
				n = len(targetFile) - (i - 1) * chunkSize
			} else {
				copy(buf, targetFile[(i - 1) * chunkSize: i * chunkSize])
				n = i * chunkSize - (i - 1) * chunkSize
			}
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
