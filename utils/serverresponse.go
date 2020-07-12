package utils

import (
	"encoding/json"
	"log"
)

type ServerResponse struct {
	Code int
	Msg  string
	Data interface{}
}

func NewServerResponse(c int, m string, d interface{}) *ServerResponse {
	return &ServerResponse{
		Code: c,
		Msg:  m,
		Data: d,
	}
}

func (sr *ServerResponse) GetInByteStream() []byte {
	bytes, err := json.Marshal(sr)
	if err != nil {
		log.Fatal("Json serialized error: ", err)
	}
	return bytes
}

func (sr *ServerResponse) GetInJsonStr() string {
	bytes, err := json.Marshal(sr)
	if err != nil {
		log.Fatal("Json serialized error: ", err)
	}
	return string(bytes)
}

func NewSimpleServerResponse(code int, msg string) *ServerResponse {
	return &ServerResponse{
		Code: code,
		Msg:  msg,
	}
}
