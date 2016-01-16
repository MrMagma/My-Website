package bloghandler

import (
    // "fmt"
    // "time"
    "database/sql"
    "encoding/json"
    "../dbhelper"
    "../goserver"
    "net/http"
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
    blogDb = blogHelper.GetDB()
    getPost, _ = blogDb.Prepare("SELECT * FROM Posts WHERE author=? AND uid=?")
)

func scanPost(rows *sql.Rows) (PostData) {
    var post PostData
    if rows.Next() {
        rows.Scan(&post.UID, &post.HTML, &post.Title, &post.Author,
            &post.Timestamp)
    }
    return post
}

func handleReq(res http.ResponseWriter, req *http.Request, args []string) {
    if (len(args) < 2) {
        return
    }
    
    postResults, _ := getPost.Query(args[0], args[1])
    
    post := scanPost(postResults)
    
    js, err := json.Marshal(post)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
    
    res.Header().Set("Content-Type", "application/json")
    res.Write(js)
}

// TODO (Joshua Gammage): Make some kind of "handler" interface so we can
// always have Register
func Register() {
    goserver.AddHandler("/blog/$/post/$", handleReq)
}
