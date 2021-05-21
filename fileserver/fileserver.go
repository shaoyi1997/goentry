package main

import (
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"github.com/valyala/fasthttp"
)

const (
	addr = "localhost:"
	port = "8080"
)

func main() {
	logger.InitLogger()

	requestHandler := generateRequestHandler()

	logger.InfoLogger.Println("HTTP file server is listening on port:", port)

	if err := fasthttp.ListenAndServe(addr+port, requestHandler); err != nil {
		logger.ErrorLogger.Fatalln("Failed in ListenAndServe:", err)
	}
}

func generateRequestHandler() func(ctx *fasthttp.RequestCtx) {
	fs := &fasthttp.FS{ //nolint:exhaustivestruct
		Root:               "./",
		GenerateIndexPages: false,
		Compress:           false,
	}

	return fs.NewRequestHandler()
}
