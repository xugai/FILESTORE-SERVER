package test

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"testing"
)

const (
	chunkSize = 1024 * 1024
	dir = "/Users/behe/Desktop/work_station/FILESTORE-SERVER/tmp/"
)

func TestByte(t *testing.T) {
	targetFile := "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/6B95694C-F72F-49F2-BFD1-D87CDE3B64F7.jpeg"
	file, err := os.Open(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	bufReader := bufio.NewReader(file)
	buf := make([]byte, chunkSize)
	k := 1
	for {
		n, err := bufReader.Read(buf)
		if err == io.EOF {
			break
		}
		newBuf := make([]byte, chunkSize)
		copy(newBuf, buf)
		fpath := dir + strconv.Itoa(k)
		targetFile, err := os.Create(fpath)
		if err != nil {
			t.Fatal(err)
		}
		targetFile.Write(newBuf[:n])
		k++
	}
}
