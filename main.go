package main

import (
	// "fmt"
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"io"
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

func testSig() {
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
				log.Printf("Process %d quited with status %d", pid, wstatus)
			}
			// all children have done
			if children == 0 {
				break
			}
		}
		done <- true
	})

	// syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)

	binary := "cat"
	path, err := exec.LookPath(binary)
	if err != nil {
		log.Fatal(err)
	}

	// var attr os.ProcAttr
	// args := []string{"sleep", "5"}
	// os.StartProcess or os/exec.Command both cannot waitpid(-1, ...)
	// proc, _ := os.StartProcess(path, args, &attr)
	// log.Printf("Process %d started", proc.Pid)
	for i:=0; i<1; i++ {
		var attr syscall.ProcAttr
		// inherit these fds , or child that need output to stdout would crash 
		attr.Files = []uintptr{0, 1, 2}
		args := []string{binary}
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


func testPipe() {
	childDone := make(chan bool, 1)
	children := 0

	binary := "sleep"
	path, err := exec.LookPath(binary)
	if err != nil {
		log.Fatal(err)
	}

	for i:=0; i<1; i++ {

		lIn, rOut, err := os.Pipe()
		if err != nil {
			log.Println(err)
		}
		var attr os.ProcAttr
		attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr, rOut}
		args := []string{binary, "2"}
		proc, err := os.StartProcess(path, args, &attr)
		if err != nil {
			log.Println(err)
		}
		// close rOut after fork, decrease it's referrence num
		rOut.Close()
		log.Printf("Process %d started", proc.Pid)
		children++

		go func(pipe *os.File){
			buf := make([]byte, 8)
			io.ReadFull(pipe, buf)
			var wstatus syscall.WaitStatus
		retry:
			pid, err := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
			if pid<=0 {
				log.Println("no exited child")
				goto retry
			}
			if err != nil {
				log.Println(err)
			}
			log.Printf("Process %d quited with status %d", pid, wstatus)
			childDone <- true
		}(lIn)
	}
	// time.Sleep(8 * time.Second)
	for i:=children; i>0; i-- {
		<-childDone
		children--
	}
}

func main() {
	// testSig()
	testPipe()
}
