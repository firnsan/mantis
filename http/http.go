package http

import (
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/service/get", serviceGetHandler)
	http.HandleFunc("/service/run", serviceRunHandler)
}

func Start() {
	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatalln(s.ListenAndServe())
}
