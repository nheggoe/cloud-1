package router

import (
	"cloud1/endpoint"
	"log"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		w.Write([]byte("Server is running!"))
	})
	router.HandleFunc("/countryinfo/v1/status", endpoint.StatusHandler)
	router.HandleFunc("/countryinfo/v1/info", endpoint.InfoHandler)
	//router.HandleFunc("/countryinfo/v1/exchange", handler)
	return router
}
