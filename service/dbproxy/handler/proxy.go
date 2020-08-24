package handler

import (
	"FILESTORE-SERVER/service/dbproxy/mapper"
	"FILESTORE-SERVER/service/dbproxy/proto"
	"bytes"
	"context"
	"encoding/json"
)

type DBProxy struct {

}

func (d *DBProxy) ExecuteAction(ctx context.Context, req *proto.ReqExec, resp *proto.RespExec) error {
	results := make([]mapper.ExecResult, len(req.Action))
	for idx, singleAction := range req.Action {
		params := []interface{}{}
		decoder := json.NewDecoder(bytes.NewReader(singleAction.Params))
		decoder.UseNumber()
		if err := decoder.Decode(&params); err != nil {
			results[idx] = mapper.ExecResult{
				Code: -2,
				Message: "请求参数有误,请重试!",
			}
			continue
		}
		// 将Number类型的参数统一转为int64，避免出现非必要的float64
		for k, param := range params {
			if v, ok := param.(json.Number); ok {
				params[k], _ = v.Int64()
			}
		}
		result, err := mapper.FuncCall(singleAction.Name, params...)
		if err != nil {
			results[idx] = mapper.ExecResult{
				Code: -2,
				Message: "dbproxy service - 服务端rpc调用失败!",
			}
			continue
		}
		results = append(results, result[0].Interface().(mapper.ExecResult))
	}
	data, _ := json.Marshal(results)
	resp.Code = 0
	resp.Message = "OK"
	resp.Data = data
	return nil
}

