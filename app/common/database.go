package common

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
    DB, _ = sql.Open("mysql", "root:@/")

    _, err := DB.Query("CREATE DATABASE IF NOT EXISTS logictree")
    CheckError(err, 3)

    _, err = DB.Query("CREATE TABLE IF NOT EXISTS logictree.conditions (field VARCHAR(255), operator VARCHAR(3), value VARCHAR(255), type VARCHAR(255), lt INT(11), rt INT(11))")
    CheckError(err, 3)

    _, err = DB.Query("CREATE TABLE IF NOT EXISTS logictree.users (name VARCHAR(255), age INT(11), num_pets INT(11))")
    CheckError(err, 3)
}