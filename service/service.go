package service

import (
	"os"
	"os/exec"
	"io"
	"syscall"
	"log"
	"errors"
	"fmt"
	"strings"
	"github.com/firnsan/fileutil"
)

func GetService(name string, gitUrl string) error {
	if name == "" || gitUrl == "" {
		return errors.New("empty name or git")
	}
	dir := "services/" + name
	if !fileutil.IsExist(dir) {
		// 下载
		cmd := exec.Command("git", "clone", gitUrl, dir)
		err := cmd.Run()
		if err != nil {
			msg := "failed to git clone: " + name
			log.Println(msg)
			return errors.New(msg)
		}
	} else {
		// return errors.New("already exists: " + name)
		// update this service
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func BuildService(name string, buildCmd string)  error {
	if name == "" || buildCmd == "" {
		return errors.New("empty name or buildCmd")
	}

	dir := "services/" + name

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		msg := "failed to build: " + name
		log.Println(msg)
		return errors.New(msg)
	}

	return nil
}

func RunService(name string, command string) (int, error) {
	if name == "" || command == "" {
		return -1, errors.New("empty name or command")
	}

	dir := "services/" + name
	args := strings.Split(command, " ")
	binary := args[0]

	if !fileutil.IsExist(dir) {
		return -1, errors.New("service not found: " + name)
	}

	lIn, rOut, err := os.Pipe()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	path, err := exec.LookPath(dir + "/" + binary)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	_ = path

	var attr os.ProcAttr
	// 切换到service所在目录
	attr.Dir = dir
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr, rOut}
	proc, err := os.StartProcess(binary, args, &attr)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	// close rOut after fork, decrease it's referrence num
	rOut.Close()
	msg := fmt.Sprintf("Process %d started", proc.Pid)
	log.Printf(msg)

	go func(){
		buf := make([]byte, 8)
		io.ReadFull(lIn, buf)
		var wstatus syscall.WaitStatus
	retry:
		pid, _ := syscall.Wait4(-1, &wstatus, syscall.WNOHANG, nil)
		if pid<=0 {
			// log.Println("no exited child")
			// need to retry because process's exit is after file closing'
			goto retry
		}
		log.Printf("Process %d quited with status %d", pid, wstatus)
	}()

	return proc.Pid, nil
}