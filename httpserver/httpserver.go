package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/config"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/routers"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/valyala/fasthttp"
)

var (
	compress    = flag.Bool("compress", false, "Whether to enable transparent response compression")
	quitChannel = make(chan os.Signal)
	rpcClient   rpc.IRPCClient
)

// TODO: shutdown routine.
func main() {
	logger.InitLogger()
	config.InitConfig()
	var err error

	rpcClient, err = rpc.NewRPCClient()
	if err != nil {
		logger.ErrorLogger.Panicln("Failed to create rpc client:", err)
	}

	router := routers.InitRouter(rpcClient)
	handler := router.Handler

	if *compress {
		handler = fasthttp.CompressHandler(handler)
	}

	go monitorForGracefulShutdown()
	logger.InfoLogger.Println("HTTP server is listening on port:", 80)
	logger.ErrorLogger.Fatalln(fasthttp.ListenAndServe(":80", handler))
}

func monitorForGracefulShutdown() {
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quitChannel
	rpcClient.Close()
	logger.InfoLogger.Println("HTTP server successfully shutdown")
	os.Exit(0)
}
