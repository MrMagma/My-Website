package goserver

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "strconv"
)

type Handler struct {
    pathRegex *regexp.Regexp
    handler func(http.ResponseWriter, *http.Request, []string)
}

var handlers = make(map[string][]Handler)

var pathStart = regexp.MustCompile("^(.+?)(?:/.*?\\$|$)")
var escapeChars = regexp.MustCompile("(\\|\\+|\\-|\\?|\\&)")

func ExtractPathData(path string) (*regexp.Regexp, string) {
    var (
        pathRegex = regexp.MustCompile("^" +
            strings.Replace(escapeChars.ReplaceAllString(path, "\\$1"), "$",
            "(.+?)", -1) + "$")
        listenMatch = pathStart.FindStringSubmatch(path)
        listenPath = "/"
    )
    
    if listenMatch != nil {
        listenPath = listenMatch[1] + "/"
    }
    
    if listenPath[len(listenPath) - 1:len(listenPath)] != "/" {
        listenPath += "/"
    }
    
    return pathRegex, listenPath
}

func CheckPath(path string) (bool) {
    _, exists := handlers[path]
    if !exists {
        handlers[path] = []Handler{}
    }
    
    return exists
}

func CallHandler(handler Handler, res http.ResponseWriter,
    req *http.Request) (bool) {
    var args = handler.pathRegex.FindStringSubmatch(req.URL.Path)
    
    if len(args) == 0 {
        return false
    }
    
    args = args[1:]
    handler.handler(res, req, args)
    
    return true
}

func InitHandler(listenPath string) {
    http.HandleFunc(listenPath, func(res http.ResponseWriter,
        req *http.Request) {
        fmt.Println("Handling request on '" + listenPath + "'...")
        for _, handler := range handlers[listenPath] {
            if CallHandler(handler, res, req) {
                return
            }
        }
    })
}

func AddHandler(path string, handler func(http.ResponseWriter,
    *http.Request, []string)) {
    
    pathRegex, listenPath := ExtractPathData(path)
    
    exists := CheckPath(listenPath)
    
    if !exists {
        InitHandler(listenPath)
    }
    
    fmt.Println("Adding handler on '" + listenPath + "'...")
    
    handlers[listenPath] = append(handlers[listenPath],
        Handler{pathRegex: pathRegex, handler: handler})
    
}

func StartServer(port int) {
    var strPort = strconv.Itoa(port)
    fmt.Println("Starting server on port " + strPort + "...")
    http.ListenAndServe(":" + strPort, nil)
}
