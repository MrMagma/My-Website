package main

import (
    // "fmt"
    // "time"
    "./dbhelper"
)

func main() {
    helper := dbhelper.NewHelper("./db/blog.db")
    helper.InitTable("Posts", []string{"uid INTEGER PRIMARY KEY",
        "html TEXT", "title TEXT", "author TEXT", "timestamp INTEGER"})
    
}
