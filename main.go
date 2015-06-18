package main


import (
	"github.com/firnsan/mantis/http"
	"log"
)

func main() {
	log.Println("start http server")
	
	go http.Start()

	select {}
}