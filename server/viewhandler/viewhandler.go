package viewhandler

import (
    "net/http"
    "path"
    "../goserver"
)

var (
    clientDir = "client"
    listenPath = "/"
)

func handleReq(res http.ResponseWriter, req *http.Request, args []string) {
    resPath := path.Join("client", req.URL.Path[len(listenPath):])
    http.ServeFile(res, req, resPath)
}

func Register() {
    goserver.AddHandler(listenPath + "$", handleReq)
}
