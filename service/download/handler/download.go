package handler

import (
	dbCli "FILESTORE-SERVER/service/dbproxy/client"
	"FILESTORE-SERVER/service/download/proto"
	"FILESTORE-SERVER/store/oss"
	"context"
	"log"
)

type Download struct {
}

const (
	objectKeyPrefix = "oss/image/"
)

func (d *Download) 	DownloadURL(ctx context.Context, req *proto.ReqDownloadURL, resp *proto.RespDownloadURL) error {
	//1. 获得用户传过来的filehash
	filehash := req.Filehash
	//2. 通过filehash去db里面查询相应的文件的path
	execResult, err := dbCli.GetFileMeta(filehash)
	if err != nil {
		log.Printf("%v\n", err)
		resp.Code = -2
		resp.Msg = "Get download url failed, please check log to get more details!"
		return err
	}
	//3. 然后通过filepath获得ali oss的signed download url
	fileMeta := dbCli.ToFileMeta(execResult.Data)
	signedURL := oss.Download(objectKeyPrefix + fileMeta.FileName)
	//4. 最后将signed url返回给用户
	resp.Code = 0
	resp.Msg = "Succeed"
	resp.Url = signedURL
	return nil
}
