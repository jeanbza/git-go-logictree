package common

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
    DB, _ = sql.Open("mysql", "root:@/")

    DB.Query("CREATE DATABASE logictree")
    DB.Query("CREATE TABLE IF NOT EXISTS logictree.conditions (field VARCHAR(255), operator VARCHAR(3), value VARCHAR(255), type VARCHAR(255), left INT(11), right INT(11)")
}