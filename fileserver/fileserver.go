package main

import (
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"github.com/valyala/fasthttp"
)

func main() {
	logger.InitLogger()
	requestHandler := generateRequestHandler()

	logger.InfoLogger.Println("Starting HTTP server on", "localhost:8080")
	if err := fasthttp.ListenAndServe("localhost:8080", requestHandler); err != nil {
		logger.ErrorLogger.Fatalln("Failed in ListenAndServe:", err)
	}
}

func generateRequestHandler() func(ctx *fasthttp.RequestCtx) {
	fs := &fasthttp.FS{
		Root:               "./images",
		GenerateIndexPages: true,
		Compress:           false,
	}

	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		default:
			fsHandler(ctx)
		}
	}

	return requestHandler
}
