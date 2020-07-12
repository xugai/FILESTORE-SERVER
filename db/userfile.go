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
