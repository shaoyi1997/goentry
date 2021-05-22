package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	pprof       = flag.Bool("pprof", false, "Whether to enable pprof")
	quitChannel = make(chan os.Signal)
	rpcClient   rpc.IRPCClient
)

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

	if *pprof {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
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
