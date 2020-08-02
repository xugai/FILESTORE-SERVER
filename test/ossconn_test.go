package test

import (
	"FILESTORE-SERVER/store/oss"
	"os"
	"testing"
)

var (
	objectKey = "oss/image/6B95694C-F72F-49F2-BFD1-D87CDE3B64F7.jpeg"
)

func TestOssconn(t *testing.T) {
	bucket := oss.Bucket()
	targetFile := "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/6B95694C-F72F-49F2-BFD1-D87CDE3B64F7.jpeg"
	file, err := os.Open(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	err = bucket.PutObject(objectKey, file)
	if err != nil {
		t.Fatal(err)
	}
}
