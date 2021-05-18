package main

import (
	"flag"

	"git.garena.com/shaoyihong/go-entry-task/httpserver/config"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/routers"

	"github.com/valyala/fasthttp"
)

var (
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	router := routers.InitRouter()
	handler := router.Handler
	if *compress {
		handler = fasthttp.CompressHandler(handler)
	}
	logger.InfoLogger.Println("HTTP server is listening at port:", 80)
	logger.ErrorLogger.Fatalln(fasthttp.ListenAndServe(":80", handler))
}
