package http

import (
	"net/http"
	_ "log"
	"fmt"
	"strconv"
	"github.com/firnsan/mantis/service"
)


func instanceRunHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}

	serviceName := req.FormValue("service")
	command := req.FormValue("cmd")
	gitUrl := req.FormValue("git")
	autoUpdate := req.FormValue("autoUpdate")
	buildCmd := req.FormValue("buildCmd")
	autoBuild := req.FormValue("autoBuild")

	if  serviceName == "" {
		http.Error(res, "empty service", http.StatusBadRequest)
		return
	}
	if  command == "" {
		http.Error(res, "empty command", http.StatusBadRequest)
		return
	}

	// 下载服务
	if gitUrl != "" {
		err := service.GetService(serviceName, gitUrl, autoUpdate=="true")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 构建服务
	if buildCmd != "" && autoBuild == "true" {
		err := service.BuildService(serviceName, buildCmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 运行服务
	pid, err := service.RunInstance(serviceName, command)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Process %d started", pid)
	res.Write([]byte(msg))
}


func instanceListHandler(res http.ResponseWriter, req *http.Request) {
	str := service.ListInstance()
	res.Write([]byte(str))
}

func instanceShutdownHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}
	pidStr := req.FormValue("pid")
	if pidStr == "" {
		http.Error(res, "empty pid", http.StatusBadRequest)
		return
	}
	pid, _ := strconv.Atoi(pidStr)
	err := service.ShutdownInstance(pid)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return

	}
	res.Write([]byte("success shutdown: " + pidStr))
}