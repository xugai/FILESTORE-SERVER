package db

import (
	"FILESTORE-SERVER/db/mysql"
	"database/sql"
	"fmt"
)

type TableFile struct {
	FileHash sql.NullString
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func OnFileUploadFinished(fileHash string, fileName string, fileSize int64, fileAddr string) bool{
	prepare, err := mysql.GetDBConnection().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`," +
			"`file_addr`, `status`) values(?, ?, ?, ?, 1)")
	if err != nil {
		fmt.Printf("Get prepare statement failed: %v\n", err)
		return false
	}
	defer prepare.Close()
	result, err := prepare.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Printf("Execute prepare statement failed: %v\n", err)
	}
	if rowsAffectedCount, err := result.RowsAffected(); err == nil {
		if rowsAffectedCount <= 0 {
			fmt.Printf("file with hash value: %v had uploaded before\n", fileHash)
		}
		return true
	}
	return false
}

func GetFileMeta(filehash string) (*TableFile, error) {
	prepare, err := mysql.GetDBConnection().Prepare("select file_sha1, file_name, file_size, file_addr" +
		" from tbl_file where file_sha1 = ? and status = 1")
	if err != nil {
		fmt.Printf("Get prepare statement failed: %v\n", err)
		return &TableFile{}, err
	}
	defer prepare.Close()
	tableFile := new(TableFile)
	err = prepare.QueryRow(filehash).Scan(&tableFile.FileHash, &tableFile.FileName, &tableFile.FileSize, &tableFile.FileAddr)
	if err != nil {
		fmt.Printf("Query file meta error: %v\n", err)
		return &TableFile{}, err
	}
	return tableFile, nil
}
