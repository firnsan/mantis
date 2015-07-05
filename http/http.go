package http

import (
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/service/get", serviceGetHandler)
	http.HandleFunc("/instance/run", instanceRunHandler)
	http.HandleFunc("/instance/list", instanceListHandler)
	http.HandleFunc("/instance/shutdown", instanceShutdownHandler)
}

func Start() {
	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatalln(s.ListenAndServe())
}
