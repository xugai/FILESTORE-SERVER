package common

//存储类型，表示文件存储到哪种媒介里
type StoreType int

const (
	_ StoreType = iota
	StoreLocal
	StoreCeph
	StoreOSS
	StoreMix
	StoreAll
)
