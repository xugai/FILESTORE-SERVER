package mapper

import "database/sql"

type TableFile struct {
	FileHash sql.NullString
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

type User struct {
	UserName string
	Email string
	Phone string
	SignupAt string
	LastActive string
	Status int
}

type UserFile struct {
	UserName string
	FileName string
	FileHash string
	FileSize int
	UploadAt string
	LastUpdate string
}

type ExecResult struct {
	Suc bool
	Code int
	Message string
	Data interface{}
}
