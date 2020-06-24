package handler

import (
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// return upload file page.
		file, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error.")
			return
		}
		io.WriteString(w, string(file))
	} else if req.Method == "POST" {
		// do upload file logic.
		file, header, err := req.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get uploaded file data, err: %v\n", err)
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			Location: "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/" + header.Filename,
			FileName: header.Filename,
			UploadAt: time.Now().Format("2006-06-02"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file in /tmp/" + header.Filename +", err: %v\n", err)
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		if err != nil {
			fmt.Printf("Failed to save data to new file, err: %v\n",  err)
		}
		http.Redirect(w, req, "/file/upload/suc", http.StatusFound)
	}
}

func UploadFileSucHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Upload finish.")
}

func GetFileMetaHandler(w http.ResponseWriter, req *http.Request) meta.FileMeta{
	req.ParseForm()
	fileHash := req.Form["fileHash"][0]
	return meta.GetFileMeta(fileHash)
}
