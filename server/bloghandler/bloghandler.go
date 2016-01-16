package bloghandler

import (
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
    
    postResults, err := db.Query(getPost, args[0], args[1])
    checkErr(err)
    
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

func editPostData(res http.ResponseWriter, req *http.Request, args []string) {
    if len(args) < 2 || req.Method != "POST" {
        return
    }
    
    req.ParseForm()
    body, exists := req.Form["body"]
    
    if exists && len(body) > 0 {
        db, err := blogHelper.OpenDB()
        checkErr(err)
        db.Exec(editPost, args[0], args[1], body[0])
        defer blogHelper.CloseDB(db)
        http.Redirect(res, req, "../" + args[1], http.StatusFound)
    }
}

func createPostData(res http.ResponseWriter, req *http.Request, args []string) {
    if len(args) < 1 || req.Method != "POST" {
        return
    }
    
    req.ParseForm()
    title, exists1 := req.Form["title"]
    body, exists2 := req.Form["body"]
    
    if exists1 && exists2 && len(title) > 0 && len(body) > 0 {
        db, err := blogHelper.OpenDB()
        checkErr(err)
        db.Exec(createPost, body[0], title[0], args[0], time.Now().Unix())
        defer blogHelper.CloseDB(db)
        http.Redirect(res, req, "../../../../..", http.StatusFound)
    }
}

// TODO (Joshua Gammage): Maybe make some kind of "handler" interface so we can
// always have Register
func Register() {
    blogHelper.InitTable("Posts", []string{"uid INTEGER PRIMARY KEY",
         "html TEXT", "title TEXT", "author TEXT", "timestamp INTEGER"})
    goserver.AddHandler("/api/blog/$/post/$", readPostData)
    goserver.AddHandler("/api/blog/$/post/$/edit", editPostData)
    goserver.AddHandler("/api/blog/$/post/create", createPostData)
}

func checkErr(err error) {
    if err != nil {
        panic("Error")
    }
}
