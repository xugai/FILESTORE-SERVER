package main

import (
	"FILESTORE-SERVER/handler"
	"fmt"
	"net/http"
)

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadFileSucHandler)
	http.HandleFunc("/file/metainfo", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.QueryFileMetasHandler)
	http.HandleFunc("/file/download", handler.FileDownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileMetaDeleteHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, err: %v\n", err)
	}
}
