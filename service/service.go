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
	"time"
	"bytes"
	"encoding/json"
	"github.com/firnsan/fileutil"
)

type Service struct {
	Git string `json:"git"`
	BuildCmd string `json:"buildCmd"`
}

type Instance struct {
	Service string `json:"service"`
	Name string `json:"name"`
	Path string `json:"path"`
	Cmd string `json:"cmd"`
}

type ProcStat struct {
	Pid int
	Name string
	Start string
	Quit string
	Args string
}

var running []ProcStat
var runningMap = make(map[int]int)
var quited []ProcStat

func GetService(name string, gitUrl string, autoUpdate bool) error {
	if name == "" || gitUrl == "" {
		return errors.New("empty name or git")
	}
	servicePath := "services/" + name
	if !fileutil.IsExist(servicePath) {
		// 下载
		cmd := exec.Command("git", "clone", gitUrl, servicePath)
		err := cmd.Run()
		if err != nil {
			msg := "failed to git clone: " + name
			log.Println(msg)
			return errors.New(msg)
		}
	} else if autoUpdate {
		// return errors.New("already exists: " + name)
		// update this service
		cmd := exec.Command("git", "pull")
		cmd.Dir = servicePath
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

	servicePath := "services/" + name

	cmd := exec.Command("sh", "-c", buildCmd)
	cmd.Dir = servicePath
	err := cmd.Run()
	if err != nil {
		msg := "failed to build: " + name
		log.Println(msg)
		return errors.New(msg)
	}

	return nil
}

func spawnProc(path string, binary string, args []string) (int, error) {
	lIn, rOut, err := os.Pipe()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	binPath, err := exec.LookPath(path + "/" + binary)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	_ = binPath
	_ = binary

	var attr os.ProcAttr
	// 切换到instance所在目录
	attr.Dir = path
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
		idx := runningMap[pid]
		_ = idx

		log.Printf("Process %d quited with status %d", pid, wstatus)
	}()

	return proc.Pid, nil
}

func RunInstance(instance Instance) (int, error) {
	serviceName := instance.Service
	name := instance.Name
	path := instance.Path
	command := instance.Cmd
	
	if serviceName == "" || command == "" {
		return -1, errors.New("empty service or command")
	}

	servicePath := "services/" + serviceName
	args := strings.Split(command, " ")
	binary := args[0]

	if !fileutil.IsExist(servicePath) {
		return -1, errors.New("service not found: " + serviceName)
	}

	if path == "" {
		// TODO:设置path为tempdir
	}

	// mkdir if path not exists
	if err := fileutil.InsureDir(path); err !=  nil {
		msg := fmt.Sprintf("fail to create path: %s : %s", path, err)
		log.Println(msg)
		return -1, errors.New(msg)
	}

	// 从servicePath复制到path
	copyCmd := "cp " + servicePath + "/* " + path
	cmd := exec.Command("sh", "-c", copyCmd)
	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf("fail to copy to: %s : %s", path, err)
		log.Println(msg)
		return -1, errors.New(msg)
	}


	pid, err := spawnProc(path, binary, args);
	if err != nil {
		return -1, err
	}

	// 簿记工作
	stat := ProcStat{
		Pid : pid,
		Name : name,
		Start  : time.Now().String(),
		Args : command,
	}
	running = append(running, stat)
	runningMap[pid] = len(running) - 1

	return pid, nil
}

func ListInstance() string {
	outBuf := new(bytes.Buffer)
	err := json.NewEncoder(outBuf).Encode(running)
	if err != nil {
		log.Println("failed to marshal json")
	}
	return outBuf.String()
}

func ShutdownInstance(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = proc.Kill()
	if err != nil {
		return err
	}
	return nil
}