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
	fileMetaArrayJsonStr, err := json.Marshal(meta.GetLastFileMetas(limitCnt))
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

