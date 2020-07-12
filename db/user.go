package db

import (
	"FILESTORE-SERVER/db/mysql"
	"fmt"
)

func UserSignup(userName string, passWord string) bool {
	prepare, err := mysql.GetDBConnection().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) " +
		"values (?, ?)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return false
	}
	defer prepare.Close()
	//todo seems prepare didn't return error msg when mysql found duplicate key error.
	result, err := prepare.Exec(userName, passWord)
	if err != nil {
		fmt.Printf("Failed to insert: %v\n", err)
		return false
	}
	if rowsAffectedCount, err := result.RowsAffected(); err == nil && rowsAffectedCount > 0 {
		return true
	} else if err == nil && rowsAffectedCount == 0 {
		fmt.Printf("Duplicate username insert!\n")
		return false
	}
	return false
}

func UserSignin(userName string, passWord string) bool {
	prepare, err := mysql.GetDBConnection().Prepare("select * from tbl_user where user_name = ? and " +
		"status = 0")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return false
	}
	defer prepare.Close()
	resultRow, err := prepare.Query(userName)
	if err != nil {
		fmt.Printf("Query username error occured: %v\n", err)
		return false
	} else if resultRow == nil {
		fmt.Printf("Username or password are incorrect!\n")
		return false
	}
	parseRows := mysql.ParseRows(resultRow)
	// []uint8 == []byte ?
	if len(parseRows) > 0 && string(parseRows[0]["user_pwd"].([]byte)) == passWord {
		return true
	}
	return false
}

func FlushUserToken(userName string, token string) bool {
	prepare, err := mysql.GetDBConnection().Prepare("replace into tbl_user_token(`user_name`, `user_token`)" +
		" values (?, ?)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
	}
	defer prepare.Close()
	_, err = prepare.Exec(userName, token)
	if err != nil {
		fmt.Printf("Replace into user token failed: %v\n", err)
		return false
	}
	return true
}

type User struct {
	UserName string
	Email string
	Phone string
	SignupAt string
	LastActive string
	Status int
}

func GetUserInfo(userName string) (User, error) {
	prepare, err := mysql.GetDBConnection().Prepare("select user_name, signup_at from tbl_user where" +
		" user_name = ? and status = 0")
	user := User{}
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return user, err
	}
	defer prepare.Close()
	err = prepare.QueryRow(userName).Scan(&user.UserName, &user.SignupAt)
	if err != nil {
		fmt.Printf("Query username error occured: %v\n", err)
		return user, err
	}
	return user, nil
}

func IfTokenIsValid(username string, token string) bool {
	prepare, err := mysql.GetDBConnection().Prepare("select count(1) from tbl_user_token where user_name = ? and user_token = ?")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return false
	}
	defer prepare.Close()
	row := prepare.QueryRow(username, token)
	var result int
	if row == nil {
		fmt.Printf("This token or username is incorrect!\n")
		return false
	}
	err = row.Scan(&result)
	if err != nil {
		fmt.Printf("Get result error: %v\n", err)
		return false
	}
	return result == 1
}
