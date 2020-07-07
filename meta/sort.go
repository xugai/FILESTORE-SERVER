package meta

import "time"

type ByUploadTime []FileMeta
const baseFormat = "2020-06-25 10:47:52"
func (b ByUploadTime) Len() int {
	return len(b)
}
//todo how to set rule to achieve asc or desc sorting?
func (b ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, b[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, b[j].UploadAt)
	return iTime.UnixNano() < jTime.UnixNano()
}

func (b ByUploadTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

