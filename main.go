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

	path, err := exec.LookPath("sleep")
	if err != nil {
		log.Fatal(err)
	}

	// var attr os.ProcAttr
	// args := []string{"sleep", "5"}
	// os.StartProcess or os/exec.Command both cannot waitpid(-1, ...)
	// proc, err := os.StartProcess(path, args, &attr)

	var attr syscall.ProcAttr
	args := []string{"sleep", "5"}
	pid, err := syscall.ForkExec(path, args, &attr)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Process %d started", proc.Pid)
	log.Printf("Process %d started", pid)
	log.Printf("Waiting child to finished")
	<-done
}

