package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	// "time"
	"log"
)


type sigHandler func(c <-chan os.Signal)

func registerSig(handler sigHandler) {

	// must be buffered channel
	c := make(chan os.Signal, 1)
	go handler(c)

	signal.Notify(c, syscall.SIGCHLD)

}

func main() {
	done := make(chan bool, 1)
	registerSig(func(c <-chan os.Signal) {
		// fmt.Println("sighandler start")
		sig := <-c
		fmt.Println(sig)
		done <- true
	})

	// syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)


	cmd := exec.Command("sleep", "5")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting child to finished")
	<-done
}

