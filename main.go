package main

import (
	"cloud1/router"
	"log"
	"net/http"
)

func main() {
	mux := router.NewRouter()
	log.Print(http.ListenAndServe(":8080", mux))
}
