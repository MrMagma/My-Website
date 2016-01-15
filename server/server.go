package main

import (
    "fmt"
    // "time"
    "./dbhelper"
    "./handlers"
    "net/http"
)
var (
    blogHelper = dbhelper.NewHelper("./db/blog.db")
    blogDb = blogHelper.GetDB()
    getPost, _ = blogDb.Prepare("SELECT * FROM Posts WHERE author=? AND uid=?")
)

func handleReq(res http.ResponseWriter, req *http.Request, args []string) {
    if (len(args) < 2) {
        return
    }
    
    rows, _ := getPost.Query(args[0], args[1])
    
    for rows.Next() {
        var (
            uid int
            html string
            title string
            author string
            timestamp int
        )
        rows.Scan(&uid, &html, &title, &author, &timestamp)
        fmt.Println("Author:", author, "Title:", title, "Timestamp:", timestamp)
    }
}

func main() {
    blogHelper.InitTable("Posts", []string{"uid INTEGER PRIMARY KEY",
        "html TEXT", "title TEXT", "author TEXT", "timestamp INTEGER"})
    handlers.AddHandler("/blog/$/post/$", handleReq)
    handlers.StartServer(8080)
}
