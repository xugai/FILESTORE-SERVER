package mapper

import (
	"FILESTORE-SERVER/service/dbproxy/conn"
	"fmt"
	"time"
)

func OnUserFileUploadFinish(userName string, fileName string, fileHash string, fileSize int64) ExecResult {
	prepare, err := conn.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`, `file_name`," +
		" `file_sha1`, `file_size`, `upload_at`, `status`) values(?, ?, ?, ?, ?, 0)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	_, err = prepare.Exec(userName, fileName, fileHash, fileSize, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Printf("Insert into tbl_user_file failed: %v\n", err)
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

func GetUserFileMetas(userName string, limit int64) ExecResult {
	prepare, err := conn.DBConn().Prepare("select file_sha1, file_size, file_name, upload_at, last_update " +
		"from tbl_user_file where user_name = ? and status = 0 limit ?")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: nil,
		}
	}
	defer prepare.Close()
	rows, err := prepare.Query(userName, limit)
	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err := rows.Scan(&userFile.FileHash, &userFile.FileSize, &userFile.FileName, &userFile.UploadAt, &userFile.LastUpdate)
		if err != nil {
			fmt.Printf("Scan to user file error: %v\n", err)
			break
		}
		userFiles = append(userFiles, userFile)
	}
	return ExecResult{
		Code: 0,
		Suc: true,
		Data: userFiles,
	}
}

func GetUserFileMeta(userName string, fileHash string) ExecResult {
	prepare, err := conn.DBConn().Prepare("select file_sha1, file_size, file_name, upload_at, last_update " +
		"from tbl_user_file where user_name = ? and file_sha1 = ?")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: nil,
		}
	}
	defer prepare.Close()
	userFile := new(UserFile)
	err = prepare.QueryRow(userName, fileHash).Scan(&userFile.UserName,
		&userFile.FileName,
		&userFile.FileHash,
		&userFile.FileSize,
		&userFile.UploadAt,
		&userFile.LastUpdate)
	if err != nil {
		fmt.Printf("Query file meta info failed: %v\n", err)
	}
	return ExecResult{
		Code: 0,
		Suc: true,
		Data: userFile,
	}
}

func UpdateUserFileMeta(userName string, fileHash string, newFileName string) ExecResult {
	prepare, err := conn.DBConn().Prepare("update tbl_user_file set file_name = ? where user_name = ?" +
		" and file_sha1 = ? and status = 0")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	exec, err := prepare.Exec(newFileName, userName, fileHash)
	if err != nil {
		fmt.Printf("Update user file meta info failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	if rowsAffected, err := exec.RowsAffected(); err == nil && rowsAffected == 1 {
		return ExecResult{
			Code: 0,
			Suc: true,
		}
	}
	fmt.Printf("Update user file meta info failed: %v\n", err)
	return ExecResult{
		Code: -2,
		Suc: false,
	}
}