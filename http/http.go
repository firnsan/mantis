package http

import (
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/service/get", serviceGetHandler)
	http.HandleFunc("/service/run", serviceRunHandler)
	http.HandleFunc("/service/list", serviceListHandler)
	http.HandleFunc("/service/shutdown", serviceShutdownHandler)
}

func Start() {
	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatalln(s.ListenAndServe())
}
