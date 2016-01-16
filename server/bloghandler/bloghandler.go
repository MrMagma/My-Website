package bloghandler

import (
    "fmt"
    "time"
    "database/sql"
    "encoding/json"
    "../dbhelper"
    "../goserver"
    "net/http"
)

var (
    createPost = "INSERT INTO Posts (html, title, author, timestamp) VALUES ($1, $2, $3, $4)"
    editPost = "UPDATE Posts SET html=$3 WHERE uid=$2 AND author=$1"
    getPost = "SELECT * FROM Posts WHERE uid=$2 AND author=$1"
)

type PostData struct {
    UID       int    `json:"uid"`
    HTML      string `json:"html"`
    Title     string `json:"title"`
    Author    string `json:"author"`
    Timestamp int    `json:"timestamp"`
}

var (
    blogHelper = dbhelper.NewHelper("./db/blog.db")
)

func scanPost(rows *sql.Rows) (PostData) {
    var post PostData
    if rows.Next() {
        rows.Scan(&post.UID, &post.HTML, &post.Title, &post.Author,
            &post.Timestamp)
    }
    return post
}

func readPostData(res http.ResponseWriter, req *http.Request, args []string) {
    if (len(args) < 2) {
        return
    }
    
    db, _ := blogHelper.OpenDB()
    
    postResults, _ := db.Query(getPost, args[0], args[1])
    
    post := scanPost(postResults)
    
    postResults.Close()
    
    js, err := json.Marshal(post)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
    
    res.Header().Set("Content-Type", "application/json")
    res.Write(js)
    
    defer blogHelper.CloseDB(db)
}

func writePostData(res http.ResponseWriter, req *http.Request, args []string) {
    if req.Method != "POST" {
        return
    }
    if len(args) == 1 {
        db, _ := blogHelper.OpenDB()
        db.Exec(createPost, req.Body, "About Bob", "Bob", time.Now().Unix())
        defer blogHelper.CloseDB(db)
    } else if len(args) == 2 {
        req.ParseForm()
        fmt.Println(req.Form)
        db, _ := blogHelper.OpenDB()
        db.Exec(editPost, args[0], args[1], "Joe")
        defer blogHelper.CloseDB(db)
    }
}

// TODO (Joshua Gammage): Make some kind of "handler" interface so we can
// always have Register
func Register() {
    blogHelper.InitTable("Posts", []string{"uid INTEGER PRIMARY KEY",
         "html TEXT", "title TEXT", "author TEXT", "timestamp INTEGER"})
    goserver.AddHandler("/api/blog/$/post/$", readPostData)
    goserver.AddHandler("/api/blog/$/post/write", writePostData)
    goserver.AddHandler("/api/blog/$/post/$/edit", writePostData)
}
