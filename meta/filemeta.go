package meta

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

// 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
