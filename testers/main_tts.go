package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "../logs/log.db")
	if err != nil {
		fmt.Println(err)
	}

	stmt, err := db.Prepare("INSERT INTO logs(tmstmp,object,state) values(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := stmt.Exec(time.Now().Unix(), "living room", 1); err != nil {
		fmt.Println(err)
	}
}
