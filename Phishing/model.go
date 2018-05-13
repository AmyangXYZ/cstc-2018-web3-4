package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "phishing:0311@tcp(mariadb:3306)/phishing?charset=utf8")
	for {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// https://github.com/go-sql-driver/mysql/issues/674
	db.SetMaxIdleConns(0)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT(10) NOT NULL AUTO_INCREMENT,
			username VARCHAR(16) NULL DEFAULT NULL,
			password VARCHAR(256) NULL DEFAULT NULL,
			PRIMARY KEY (id)
		);`)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS admin (
			username VARCHAR(16) NULL DEFAULT NULL,
			password VARCHAR(256) NULL DEFAULT NULL);`)
	db.Exec(`INSERT admin SET username='admin', password='` + os.Getenv("flag") + `'`)
}

func insert(user, passwd string) {
	dangerousSQL := fmt.Sprintf(`INSERT INTO users(username,password) VALUES('%s', '%s')`, user, passwd)
	db.Exec(dangerousSQL)
}
