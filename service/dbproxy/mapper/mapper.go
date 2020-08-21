package mapper

import (
	"errors"
	"reflect"
)

var funcs = map[string]interface{}{
	"/user/UserSignup": UserSignup,
	"/user/UserSignin": UserSignin,
	"/user/FlushUserToken": FlushUserToken,
	"/user/GetUserInfo": GetUserInfo,
	"/user/IfTokenIsValid": IfTokenIsValid,

	"/file/OnFileUploadFinished": OnFileUploadFinished,
	"/file/GetFileMeta": GetFileMeta,
	"/file/UpdateFileStoreLocation": UpdateFileStoreLocation,

	"/ufile/OnUserFileUploadFinish": OnUserFileUploadFinish,
	"/ufile/GetUserFileMetas": GetUserFileMetas,
	"/ufile/GetUserFileMeta": GetUserFileMeta,
	"/ufile/UpdateUserFileMeta": UpdateUserFileMeta,
}

// 通过反射动态调用指定函数
func FuncCall(name string, params ... interface{}) ([]reflect.Value, error) {
	if _, ok := funcs[name]; !ok {
		err := errors.New("指定调用的函数不存在!")
		return nil, err
	}
	f := reflect.ValueOf(funcs[name])
	if len(params) != f.Type().NumIn() {
		err := errors.New("传入参数与调用函数的参数长度不一致!")
		return nil, err
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in), nil
}
