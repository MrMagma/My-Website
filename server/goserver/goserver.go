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

var handlers = Handlers{}

var escapeChars = regexp.MustCompile("(\\|\\+|\\-|\\?|\\&)")

func GetPathRegex(path string) (*regexp.Regexp) {
    var pathRegex = regexp.MustCompile("^" +
        strings.Replace(escapeChars.ReplaceAllString(path, "\\$1"), "$",
            "(.+?)", -1) + "$")
    
    return pathRegex
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

func AddHandler(path string, handler func(http.ResponseWriter,
    *http.Request, []string)) {
    
    pathRegex := GetPathRegex(path)
    
    fmt.Println("Adding handler on '" + path + "'...")
    
    handlers = append(handlers,
        Handler{pathRegex: pathRegex, handler: handler, path: path})
    
    sort.Sort(handlers)    
}

func handleReq(res http.ResponseWriter,
    req *http.Request) {
    for _, handler := range handlers {
        if CallHandler(handler, res, req) {
            fmt.Println("Handled request with '" + handler.path + "' handler")
            return
        }
    }
}

func Start(port int) {
    var strPort = strconv.Itoa(port)
    fmt.Println("Starting server on port " + strPort + "...")
    http.HandleFunc("/", handleReq)
    http.ListenAndServe(":" + strPort, nil)
}
