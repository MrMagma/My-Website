package goserver

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "strconv"
    "sort"
)

type Handler struct {
    path string
    pathRegex *regexp.Regexp
    handler func(http.ResponseWriter, *http.Request, []string)
}

type Handlers []Handler

func (slice Handlers) Len() int {
    return len(slice)
}

func (slice Handlers) Less(i, j int) bool {
    return len(slice[i].path) > len(slice[j].path)
}

func (slice Handlers) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

var handlers = make(map[string]Handlers)

var pathStart = regexp.MustCompile("^(.+?)(?:(?:/|).*?\\$|$)")
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
        listenPath = listenMatch[1]
    }
    
    if listenPath[len(listenPath) - 1:len(listenPath)] != "/" {
        listenPath += "/"
    }
    
    return pathRegex, listenPath
}

func CheckPath(path string) (bool) {
    _, exists := handlers[path]
    if !exists {
        handlers[path] = Handlers{}
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
                fmt.Println("Handled request with '" + handler.path + "' handler")
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
        Handler{pathRegex: pathRegex, handler: handler, path: path})
    
    sort.Sort(handlers[listenPath])    
}

func Start(port int) {
    var strPort = strconv.Itoa(port)
    fmt.Println("Starting server on port " + strPort + "...")
    http.ListenAndServe(":" + strPort, nil)
}
