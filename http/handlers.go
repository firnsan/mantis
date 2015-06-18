package http

import (
	"net/http"
	"github.com/firnsan/mantis/service"
	_ "log"
	"fmt"
)

func serviceGetHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}

	name := req.FormValue("service")
	gitUrl := req.FormValue("git")
	buildCmd := req.FormValue("build")

	if  name == "" {
		http.Error(res, "empty service", http.StatusBadRequest)
		return
	}
	if  gitUrl == "" {
		http.Error(res, "empty gitUrl", http.StatusBadRequest)
		return
	}

	// 下载服务
	err := service.GetService(name, gitUrl)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// 构建服务
	if buildCmd != "" {
		err := service.BuildService(name, buildCmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	res.Write([]byte("success get: " + name))
}

func serviceRunHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}

	name := req.FormValue("service")
	command := req.FormValue("cmd")
	gitUrl := req.FormValue("git")
	buildCmd := req.FormValue("build")

	if  name == "" {
		http.Error(res, "empty service", http.StatusBadRequest)
		return
	}
	if  command == "" {
		http.Error(res, "empty command", http.StatusBadRequest)
		return
	}

	// 下载服务
	if gitUrl != "" {
		err := service.GetService(name, gitUrl)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 构建服务
	if buildCmd != "" {
		err := service.BuildService(name, buildCmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 运行服务
	pid, err := service.RunService(name, command)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Process %d started", pid)
	res.Write([]byte(msg))
}