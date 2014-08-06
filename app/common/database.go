package common

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
    DB, _ = sql.Open("mysql", "root:@/")

    DB.Query("CREATE DATABASE logictree")
    DB.Query("CREATE TABLE IF NOT EXISTS logictree.equality (field VARCHAR(255), operator VARCHAR(3), value VARCHAR(255)")
    DB.Query("CREATE TABLE IF NOT EXISTS logictree.logic (operator VARCHAR(3))")
}