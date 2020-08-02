package test

import (
	"FILESTORE-SERVER/utils"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestJsonIterator(t *testing.T) {
	chkArr := []int{1, 2, 3}
	body := utils.NewServerResponse(200, "OK", chkArr).GetInByteStream()
	get := jsoniter.Get(body, "Data")
	arr := get.GetInterface().([]interface{})
	fmt.Println(len(arr))
	for i := 0; i < len(arr); i++ {
		fmt.Println(int(arr[i].(float64)))
	}
}
