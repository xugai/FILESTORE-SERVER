package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request from external...")
	fmt.Fprintf(w, "Hello World Berio Xu!")
}


func main() {
	log.Println("Hello World Service Started...")
	http.HandleFunc("/svc/hello", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 编译、打包镜像、编写k8s service的yaml文件、运行到k8s里面、通过外部访问
