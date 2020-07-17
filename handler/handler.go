package handler

import (
	"FILESTORE-SERVER/db"
	"FILESTORE-SERVER/meta"
	"FILESTORE-SERVER/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
		req.ParseForm()
		file, header, err := req.FormFile("file")
		userName := req.Form.Get("username")
		if err != nil {
			fmt.Printf("Failed to get uploaded file data, err: %v\n", err)
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			Location: "/Users/behe/Desktop/work_station/FILESTORE-SERVER/file/" + header.Filename,
			FileName: header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
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
		log.Printf("%v Upload file with hash: %v", time.Now().Format("2006-01-02 15:04:05"), fileMeta.FileSha1)
		//meta.UpdateFileMeta(fileMeta)
		result := meta.UpdateFileMetaDB(fileMeta)
		if !result {
			fmt.Printf("Update or insert file meta in DB error, please check.\n")
			return
		}
		if err != nil {
			fmt.Printf("Failed to save data to new file, err: %v\n",  err)
			return
		}
		result = db.OnUserFileUploadFinish(userName, fileMeta.FileName, fileMeta.FileSha1, fileMeta.FileSize)
		if result {
			http.Redirect(w, req, "/static/view/home.html", http.StatusFound)
		} else {
			fmt.Println("Upload user file error, please check log.")
		}
	}
}

func UploadFileSucHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Upload finish.")
}

func GetFileMetaHandler(w http.ResponseWriter, req *http.Request){
	req.ParseForm()
	filehash := req.Form.Get("filehash")
	//fileMeta := meta.GetFileMeta(filehash)
	fileMeta, err2 := meta.GetFileMetaDB(filehash)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Get file meta error occured: %v\n", err2)
		return
	}
	fileMetaJsonStr, err := json.Marshal(*fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(fileMetaJsonStr)
}

func QueryFileMetasHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	limitCnt, err := strconv.Atoi(req.Form.Get("limit"))
	if err != nil {
		w.Write([]byte("param: limit is invalid!"))
		return
	}
	userName := req.Form.Get("username")
	userFiles, _ := db.GetUserFileMetas(userName, limitCnt)
	fileMetaArrayJsonStr, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(fileMetaArrayJsonStr)
}

func FileDownloadHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fileHash := req.Form.Get("filehash")
	fm := meta.GetFileMeta(fileHash)
	file, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\"" + fm.FileName + "\"")
	w.Write(bytes)
}

func FileMetaUpdateHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	opType := req.Form.Get("opType")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	filehash := req.Form.Get("filehash")
	newFileName := req.Form.Get("newFileName")
	currFileMeta := meta.GetFileMeta(filehash)
	currFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currFileMeta)

	currFileMetaJsonStr, err := json.Marshal(currFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(currFileMetaJsonStr)
}

func FileMetaDeleteHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	req.ParseForm()
	filehash := req.Form.Get("filehash")
	err := os.Remove(meta.GetFileMeta(filehash).Location)
	if err != nil {
		log.Printf("Delete file err: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meta.RemoveFileMeta(filehash)
	w.WriteHeader(http.StatusOK)
}

// 尝试秒传处理用户上传的文件
func TryFastUploadHandler(w http.ResponseWriter, req *http.Request) {
	// 获取页面传回来的用户名、当前上传文件的SHA1值
	req.ParseForm()
	userName := req.Form.Get("username")
	fileHash := req.Form.Get("filehash")
	fileName := req.Form.Get("filename")
	fileSize, _ := strconv.Atoi(req.Form.Get("filesize"))
	// 根据SHA1值在tbl_file中查找是否已经存在，如果存在，则往tbl_user_file中插入新的记录，否则秒传失败
	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		fmt.Printf("Get file meta error: %v\n", err)
		serverResponse := utils.NewSimpleServerResponse(500, "文件秒传上传失败，当前上传的文件从未有人上传过!")
		w.Write(serverResponse.GetInByteStream())
		return
	}
	if fileMeta.FileHash.String == "" {
		serverResponse := utils.NewSimpleServerResponse(500, "秒传失败，请尝试普通文件上传接口")
		w.Write(serverResponse.GetInByteStream())
		return
	}
	result := db.OnUserFileUploadFinish(userName, fileName, fileMeta.FileHash.String, int64(fileSize))
	if !result {
		serverResponse := utils.NewSimpleServerResponse(500, "文件秒传上传失败，请稍后重试")
		w.Write(serverResponse.GetInByteStream())
		return
	}
	w.Write(utils.NewSimpleServerResponse(200, "文件秒传上传成功!").GetInByteStream())
	w.WriteHeader(http.StatusOK)
}

