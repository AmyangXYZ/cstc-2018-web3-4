package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AmyangXYZ/sweetygo"
	"github.com/AmyangXYZ/sweetygo/middlewares"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:0311@tcp(localhost:3306)/cswsh?charset=utf8")

	for {
		err := db.Ping()
		if err == nil {
			fmt.Println("db ok")
			break
		}
		time.Sleep(2 * time.Second)
	}

	// https://github.com/go-sql-driver/mysql/issues/674
	db.SetMaxIdleConns(0)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INT(10) NOT NULL AUTO_INCREMENT,
			msg VARCHAR(1024) NULL DEFAULT NULL,
			PRIMARY KEY (id)
		);`)
}

func main() {
	stepping := sweetygo.New("./", nil)
	target := sweetygo.New("./", nil)

	stepping.USE(middlewares.Logger(os.Stdout))
	stepping.GET("/", steppingIndex)
	stepping.GET("/static/*files", static)
	stepping.GET("/admin", adminPage)
	stepping.POST("/api/msg", saveEvilsJS)
	stepping.GET("/api/msg", readEvilJS)
	go stepping.RunServer(":8001")

	target.USE(middlewares.Logger(os.Stdout))
	target.GET("/", targetIndex)
	target.GET("/ws", ws)
	target.GET("/static/*files", static)
	target.RunServer("127.0.0.1:8002")
}

func steppingIndex(ctx *sweetygo.Context) {
	ctx.Render(200, "stepping-index")
}

func saveEvilsJS(ctx *sweetygo.Context) {
	evilJS := ctx.Param("msg")
	stmt, err := db.Prepare("INSERT comments SET msg=?")

	_, err = stmt.Exec(evilJS)
	if err != nil {
		fmt.Println(err)
	}
	ctx.JSON(201, 1, "success", nil)
}

func readEvilJS(ctx *sweetygo.Context) {
	if ctx.GetCookie("token") != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTI2MDIxNDU2LCJuYW1lIjoiQW15YW5nIn0.JJmKn7DuM1VbriXeG4XqT18ycDdObdaE1fltp2CIGAY" {
		ctx.JSON(401, 0, "unauthorized", nil)
		return
	}
	rows, err := db.Query("select msg from comments")
	var msgs []string
	var msg string
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&msg)
		msgs = append(msgs, msg)
	}
	ctx.JSON(200, 1, "success", map[string][]string{
		"msgs": msgs,
	})
}

func adminPage(ctx *sweetygo.Context) {
	if ctx.GetCookie("token") != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTI2MDIxNDU2LCJuYW1lIjoiQW15YW5nIn0.JJmKn7DuM1VbriXeG4XqT18ycDdObdaE1fltp2CIGAY" {
		ctx.JSON(401, 0, "unauthorized", nil)
		return
	}
	ctx.Render(200, "stepping-admin")
}

func targetIndex(ctx *sweetygo.Context) {
	ctx.Render(200, "target-index")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

type msg struct {
	Cmd string
}

func ws(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
			break
		}

		res := exec(m.Cmd)
		fmt.Println(res)
		if err = conn.WriteJSON(res); err != nil {
			fmt.Println(err)
			break
		}
	}
}

// false exec :)
func exec(cmd string) string {
	switch cmd {
	case "ls":
		return "flaaaaag.txt index.html jquery-3.1.1.min.js server.go"
	case "cat flaaaaag.txt":
		return "flag{ed3e359a08bdec24a1b9ed1fdb2fd90e}"
	default:
		return "Hacker Denied!"
	}
}
