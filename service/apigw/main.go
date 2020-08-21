package main

import "FILESTORE-SERVER/service/apigw/route"

func main() {
	router := route.Router()
	router.Run(":8080")
}
