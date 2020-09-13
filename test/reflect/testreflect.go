package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type FileMeta struct {
	FileName string
	FileSize int
	FileHash string
	FileLocation string
}

func main() {
	fileMeta := FileMeta{
		FileName: "a.txt",
		FileSize: 256,
		FileHash: "abcdefg",
		FileLocation: "/aaa/bbb/ccc/d",
	}
	typeOfFileMeta := reflect.TypeOf(fileMeta)
	fmt.Println("type is: ", typeOfFileMeta)
	valueOfFileMeta := reflect.ValueOf(fileMeta)
	kind := valueOfFileMeta.Kind()
	t := valueOfFileMeta.Type()
	zeroValue := reflect.Zero(t)
	isValid := valueOfFileMeta.IsValid()
	fmt.Println("value is: ", valueOfFileMeta)
	fmt.Println("kind is: ", kind)
	fmt.Println("t is: ", t)
	fmt.Println("zero value is: ", zeroValue)
	fmt.Println("is valid? ",  isValid)

	var input interface{}
	input = map[string]interface{}{
		"FileName": map[string]interface{}{"String": "a.txt"},
		"FileHash": map[string]interface{}{"String": "abcdefg"},
		"FileSize": map[string]interface{}{"Integer": 256},
		"FileLocation": map[string]interface{}{"String": "/aaa/bbb/ccc/dd"},
	}
	output := FileMeta{}
	mapstructure.Decode(input, &output)
	fmt.Println(output)
}