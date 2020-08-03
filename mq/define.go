package mq

import "FILESTORE-SERVER/common"

type TransferData struct {
	FileHash string
	CurLocation string
	DestLocation string
	DestStoreType common.StoreType
}
