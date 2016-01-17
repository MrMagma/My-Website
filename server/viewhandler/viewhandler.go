package viewhandler

import (
    "net/http"
    "path"
    "os"
    "strings"
    "errors"
    "../goserver"
    "fmt"
)

var (
    clientDir = "client"
    listenPath = "/"
)

func serve404(res http.ResponseWriter, req *http.Request) {
    resBytes, err := fetchFile(path.Join(clientDir, "404.html"))
    if err != nil {
        res.Write([]byte("404: File not found"))
    }
    res.Write(resBytes)
}

func fetchFile(resPath string) ([]byte, error) {
    file, err := os.Open(resPath)
    if err != nil {
        // If there was no extension then try one last time with an extension
        if len(path.Ext(resPath)) == 0 {
            return fetchFile(resPath + ".html")
        }
        return nil, err
    }
    
    info, err := file.Stat()
    if err != nil {
        defer file.Close()
        return nil, err
    }
    
    if info.IsDir() {
        files, err := file.Readdirnames(128)
        defer file.Close()
        
        if err != nil {
            return nil, err
        }
        
        for _, fName := range files {
            if strings.HasPrefix(fName, "index") {
                return fetchFile(path.Join(resPath, fName))
            }
        } 
        
        return nil, errors.New("index file not found in directory")
    }
    
    fBytes := make([]byte, 65536)
    _, err = file.Read(fBytes)
    defer file.Close()
    
    return fBytes, err
}

func serveClient(res http.ResponseWriter, req *http.Request, args []string) {
    resPath := path.Join("client", req.URL.Path[len(listenPath):])
    
    resBytes, err := fetchFile(resPath)
    if err != nil {
        handleErr(err)
        serve404(res, req)
        return
    }
    
    res.Write(resBytes)
}

func handleErr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

func Register() {
    goserver.AddHandler(listenPath + "$", serveClient)
}
