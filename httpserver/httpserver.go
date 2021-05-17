package main

import (
	"flag"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/routers"

	"github.com/valyala/fasthttp"
)

var (
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	logger.InitLogger()
	router := routers.InitRouter()
	handler := router.Handler
	if *compress {
		handler = fasthttp.CompressHandler(handler)
	}
	logger.InfoLogger.Println("HTTP server is listening at port:", 80)
	logger.ErrorLogger.Fatalln(fasthttp.ListenAndServe(":80", handler))
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	//
	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	//var c fasthttp.Cookie
	//c.SetKey("cookie-name")
	//c.SetValue("cookie-value")
	//ctx.Response.Header.SetCookie(&c)
}
