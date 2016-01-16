package dbhelper

import (
    "fmt"
    "database/sql"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

type DBHelper struct {
    path string
    open bool
}

func NewHelper(path string) (*DBHelper) {
    helper := new(DBHelper)
    helper.path = path
    helper.open = false
    return helper
}

func (self *DBHelper) InitTable(name string, args []string) (sql.Result,
    error) {
    db, _ := self.OpenDB()
    res, err := db.Exec("CREATE TABLE IF NOT EXISTS " + name + " (" +
        strings.Join(args[:], ", ") + ")")
    checkErr(err)
    self.CloseDB(db)
    return res, err
}

func (self *DBHelper) OpenDB() (*sql.DB, error) {
    if (self.open) {
        fmt.Println("DB already open!")
    }
    self.open = true
    db, err := sql.Open("sqlite3", self.path)
    checkErr(err)
    return db, err
}

func (self *DBHelper) CloseDB(db *sql.DB) {
    db.Close()
    self.open = false
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
