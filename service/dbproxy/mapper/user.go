package mapper

import (
	"FILESTORE-SERVER/service/dbproxy/conn"
	"fmt"
)

func UserSignup(userName string, passWord string) ExecResult {
	prepare, err := conn.DBConn().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) " +
		"values (?, ?)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	//todo seems prepare didn't return error msg when mysql found duplicate key error
	//todo because you use 'insert ignore' syntax, it will ignore error when mysql execute command
	result, err := prepare.Exec(userName, passWord)
	if err != nil {
		fmt.Printf("Failed to insert: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	if rowsAffectedCount, err := result.RowsAffected(); err == nil && rowsAffectedCount > 0 {
		return ExecResult{
			Code: 0,
			Suc: true,
		}
	} else if err == nil && rowsAffectedCount == 0 {
		fmt.Printf("Duplicate username insert!\n")
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	return ExecResult{
		Code: -2,
		Suc: false,
	}
}

func UserSignin(userName string, passWord string) ExecResult {
	prepare, err := conn.DBConn().Prepare("select * from tbl_user where user_name = ? and " +
		"status = 0")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	resultRow, err := prepare.Query(userName)
	if err != nil {
		fmt.Printf("Query username error occured: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	} else if resultRow == nil {
		fmt.Printf("Username or password are incorrect!\n")
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	parseRows := conn.ParseRows(resultRow)
	// []uint8 == []byte ?
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == passWord {
		return ExecResult{
			Code: 0,
			Suc: true,
		}
	}
	return ExecResult{
		Code: -2,
		Suc: false,
	}
}

func FlushUserToken(userName string, token string) ExecResult {
	prepare, err := conn.DBConn().Prepare("replace into tbl_user_token(`user_name`, `user_token`)" +
		" values (?, ?)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
	}
	defer prepare.Close()
	_, err = prepare.Exec(userName, token)
	if err != nil {
		fmt.Printf("Replace into user token failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	return ExecResult{
		Code: 0,
		Suc: true,
	}
}

func GetUserInfo(userName string) ExecResult {
	prepare, err := conn.DBConn().Prepare("select user_name, signup_at from tbl_user where" +
		" user_name = ? and status = 0")
	user := User{}
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: user,
		}
	}
	defer prepare.Close()
	err = prepare.QueryRow(userName).Scan(&user.UserName, &user.SignupAt)
	if err != nil {
		fmt.Printf("Query username error occured: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: user,
		}
	}
	return ExecResult{
		Code: 0,
		Suc: true,
		Data: user,
	}
}

func IfTokenIsValid(username string, token string) ExecResult {
	prepare, err := conn.DBConn().Prepare("select count(1) from tbl_user_token where user_name = ? and user_token = ?")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	row := prepare.QueryRow(username, token)
	var result int
	if row == nil {
		fmt.Printf("This token or username is incorrect!\n")
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	err = row.Scan(&result)
	if err != nil {
		fmt.Printf("Get result error: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	if result == 1 {
		return ExecResult{
			Code: 0,
			Suc: true,
		}
	}
	return ExecResult{
		Code: -2,
		Suc: false,
	}
}
