package db

import (
	"FILESTORE-SERVER/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName string
	FileName string
	FileHash string
	FileSize int
	UploadAt string
	LastUpdate string
}

func OnUserFileUploadFinish(userName string, fileName string, fileHash string, fileSize int64) bool {
	prepare, err := mysql.GetDBConnection().Prepare("insert ignore into tbl_user_file (`user_name`, `file_name`," +
		" `file_sha1`, `file_size`, `upload_at`, `status`) values(?, ?, ?, ?, ?, 0)")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return false
	}
	defer prepare.Close()
	_, err = prepare.Exec(userName, fileName, fileHash, fileSize, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Printf("Insert into tbl_user_file failed: %v\n", err)
		return false
	}
	return true
}

func GetUserFileMetas(userName string, limit int) ([]UserFile, error) {
	prepare, err := mysql.GetDBConnection().Prepare("select file_sha1, file_size, file_name, upload_at, last_update " +
		"from tbl_user_file where user_name = ? and status = 0 limit ?")
	if err != nil {
		fmt.Printf("Prepare statement failed: %v\n", err)
		return nil, err
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
	return userFiles, nil
}