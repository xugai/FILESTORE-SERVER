package handler

import (
	"FILESTORE-SERVER/service/account/proto"
	dbCli "FILESTORE-SERVER/service/dbproxy/client"
	"context"
	"encoding/json"
	"log"
)

// todo 问题记录：定义的rpc sever API，因为在proto文件里定义的返回值是一个error，但如果在实现体里面的返回值是以errors为首的新new出来的error，用作返回值返回的话，好像相应的resp是nil的
// todo 问题更新： 如果rpc server API 在执行过程中返回的值是非nil的error，那么go micro fwk好像不会对resp结构体进行封装，或者说，此时rpc client API返回的resp是nil的
func (u *User) QueryFileMetas(ctx context.Context, req *proto.ReqQueryFileMetas, resp *proto.RespQueryFileMetas) error {
	username := req.Username
	limitCnt := req.LimitCnt
	execResult, err := dbCli.GetUserFileMetas(username, int(limitCnt))
	if err != nil {
		log.Println(err)
		resp.Code = -2
		resp.Msg = "服务内部发生错误，请检查错误日志"
		return err
	}
	fileMetas, _ := json.Marshal(dbCli.ToTableUserFiles(execResult.Data))
	resp.Code = 0
	resp.Msg = "Succeed"
	resp.FileMetas = fileMetas
	return nil
}
