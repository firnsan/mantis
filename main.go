package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"
)


type sigHandler func(c <-chan os.Signal)

func registerSig(handler sigHandler) {
	fmt.Println("in registerSig")

	// must be buffered channel
	c := make(chan os.Signal, 1)
	go handler(c)

	signal.Notify(c, syscall.SIGCHLD)

}

func main() {
	registerSig(func(c <-chan os.Signal) {
		// fmt.Println("sighandler start")
		sig := <-c
		fmt.Println(sig)
	})

	syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)

	time.Sleep(2 * time.Second)
}

