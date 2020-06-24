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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, err: %v\n", err)
	}
}
