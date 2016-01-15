package dbhelper

import (
    // "fmt"
    // "time"
    "database/sql"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

type DBHelper struct {
    db *sql.DB
}

func NewHelper(path string) (*DBHelper) {
    var err error
    helper := new(DBHelper)
    helper.db, err = OpenDB(path)
    checkErr(err)
    return helper
}

func (self *DBHelper) GetDB() (*sql.DB) {
    return self.db
}

func (self *DBHelper) InitTable(name string, args []string) (res sql.Result,
    err error) {
    res, err = self.db.Exec("CREATE TABLE IF NOT EXISTS " + name + " (" +
        strings.Join(args[:], ", ") + ")")
    return res, err
}

func OpenDB(path string) (db *sql.DB, err error) {
    db, err = sql.Open("sqlite3", path)
    return db, err
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
