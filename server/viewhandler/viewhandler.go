package viewhandler

import (
    "fmt"
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
    fmt.Println(resPath)
}

func Register() {
    goserver.AddHandler(listenPath, handleReq)
    goserver.AddHandler(listenPath + "$", handleReq)
}
