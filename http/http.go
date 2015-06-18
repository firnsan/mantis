package http

import (
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/run", runHandler)
}

func Start() {
	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatalln(s.ListenAndServe())
}
