package http

import (
	"net/http"
	_ "log"
	_ "fmt"
	_ "strconv"
	"github.com/firnsan/mantis/service"
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
	err := service.GetService(name, gitUrl, true)
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