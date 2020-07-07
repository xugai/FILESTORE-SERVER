package meta

import (
	"FILESTORE-SERVER/db"
	"sort"
)

type FileMeta struct {
	UploadAt string
	FileName string
	FileSha1 string
	FileSize int64
	Location string
}

var fileMetas = map[string]FileMeta{}

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增、更新文件的元信息
func UpdateFileMeta(fileMeta FileMeta) {
	fileMetas[fileMeta.FileSha1] = fileMeta
}

func UpdateFileMetaDB(fileMeta FileMeta) bool{
	return db.OnFileUploadFinished(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location)
}

// 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tableFile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return &FileMeta{}, err
	}
	fileMeta := &FileMeta{
		FileSha1: tableFile.FileHash.String,
		FileName: tableFile.FileName.String,
		FileSize: tableFile.FileSize.Int64,
		Location: tableFile.FileAddr.String,
	}
	return fileMeta, nil
}

func GetLastFileMetas(limitCnt int) []FileMeta {
	fileMetaArray := make([]FileMeta, 0)
	for _, v := range fileMetas {
		fileMetaArray = append(fileMetaArray, v)
	}
	sort.Sort(ByUploadTime(fileMetaArray))
	if limitCnt >= len(fileMetaArray) {
		return fileMetaArray
	}
	return fileMetaArray[0:limitCnt]
}

func RemoveFileMeta(filehash string) {
	delete(fileMetas, filehash)
}


