package mapper

import (
	"FILESTORE-SERVER/service/dbproxy/conn"
	"fmt"
)

func OnFileUploadFinished(fileHash string, fileName string, fileSize int64, fileAddr string) ExecResult {
	prepare, err := conn.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`," +
			"`file_addr`, `status`) values(?, ?, ?, ?, 1)")
	if err != nil {
		fmt.Printf("Get prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
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

func GetFileMeta(filehash string) ExecResult {
	prepare, err := conn.DBConn().Prepare("select file_sha1, file_name, file_size, file_addr" +
		" from tbl_file where file_sha1 = ? and status = 1")
	if err != nil {
		fmt.Printf("Get prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: TableFile{},
		}
	}
	defer prepare.Close()
	tableFile := new(TableFile)
	err = prepare.QueryRow(filehash).Scan(&tableFile.FileHash, &tableFile.FileName, &tableFile.FileSize, &tableFile.FileAddr)
	if err != nil {
		fmt.Printf("Query file meta error: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
			Data: TableFile{},
		}
	}
	return ExecResult{
		Code: 0,
		Suc: true,
		Data: tableFile,
	}
}

func UpdateFileStoreLocation(fileHash, fileAddr string) ExecResult {
	prepare, err := conn.DBConn().Prepare("update tbl_file set file_addr = ? where file_sha1 = ?")
	if err != nil {
		fmt.Printf("Get prepare statement failed: %v\n", err)
		return ExecResult{
			Code: -2,
			Suc: false,
		}
	}
	defer prepare.Close()
	_, err = prepare.Exec(fileAddr, fileHash)
	if err != nil {
		fmt.Printf("Execute prepare statement failed: %v\n", err)
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
