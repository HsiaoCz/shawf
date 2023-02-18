package main

import (
	"net/http"
	"shawf/shawg/shawg"
)

func main() {
	r := shawg.New()
	r.GET("/", func(ctx *shawg.Context) {
		ctx.String(http.StatusOK, "<h1>Hello My man</h1>")
	})
	r.GET("/hello", func(ctx *shawg.Context) {
		ctx.JSON(http.StatusOK, shawg.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})
	r.GET("/hi/:name", func(ctx *shawg.Context) {
		ctx.JSON(http.StatusOK, shawg.H{
			"code": 100,
			"msg":  "Hello:" + ctx.Param("name"),
		})
	})
	r.Run(":3023")
}
