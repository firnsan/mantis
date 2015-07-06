package http

import (
	"net/http"
	_ "log"
	_ "fmt"
	_ "strconv"
	"encoding/json"
	"github.com/firnsan/mantis/service"
)

func serviceDeployHandler(res http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(res, "empty body", http.StatusBadRequest)
		return
	}

	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)

	serv := new(service.Service)
	json.Unmarshal(buf, serv)


	name := serv.Name
	gitUrl := serv.Git
	buildCmd := serv.BuildCmd

	if  name == "" {
		http.Error(res, "empty service", http.StatusBadRequest)
		return
	}
	if  gitUrl == "" {
		http.Error(res, "empty gitUrl", http.StatusBadRequest)
		return
	}

	// 下载服务
	err := service.DeployService(name, gitUrl, true)
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