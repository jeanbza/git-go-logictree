package common

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
    DB, _ = sql.Open("mysql", "root:@/")

    DB.Query("CREATE DATABASE logictree")
    DB.Query("CREATE TABLE IF NOT EXISTS logictree.equality (field VARCHAR(255), operator VARCHAR(5), value FLOAT(25))")
    DB.Query("CREATE TABLE IF NOT EXISTS logictree.logic (operator VARCHAR(3))")
}