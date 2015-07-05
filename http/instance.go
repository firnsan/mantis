package http

import (
	"net/http"
	_ "log"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/firnsan/mantis/service"
)


func instanceRunHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}

	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)

	instance := new(service.Instance)
	json.Unmarshal(buf, instance)


	if  instance.Service == "" {
		http.Error(res, "empty service", http.StatusBadRequest)
		return
	}
	if  instance.Cmd == "" {
		http.Error(res, "empty command", http.StatusBadRequest)
		return
	}

	/*
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
*/
	// 运行服务
	pid, err := service.RunInstance(instance)
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