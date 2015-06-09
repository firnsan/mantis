package main

import (
	// "fmt"
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
	children := 0
	registerSig(func(c <-chan os.Signal) {
		// func to handler SIGCHLD
		for {
			_ = <-c

			for {
				var wstatus syscall.WaitStatus
				pid, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
				if pid<=0 {
					break
				}
				if err != nil {
					log.Println(err)
					break
				}
				children--
				log.Printf("Process %d quited", pid)
			}
			// all children have done
			if children == 0 {
				break
			}
		}
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
	// proc, _ := os.StartProcess(path, args, &attr)
	// log.Printf("Process %d started", proc.Pid)
	for i:=0; i<8; i++ {
		var attr syscall.ProcAttr
		args := []string{"sleep", "5"}
		pid, err := syscall.ForkExec(path, args, &attr)
		if err != nil {
			log.Fatal(err)
		}
		children++
		log.Printf("Process %d started", pid)
	}
	log.Printf("Waiting child to finished")
	<-done
}

