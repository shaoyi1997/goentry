package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services"
)

var (
	quitChannel         = make(chan os.Signal)
	waitGroup           = sync.WaitGroup{}
	isShutdownInitiated = false
)

func main() {
	config.InitConfig()
	services.Init()
	initTCPServer()
}

func initTCPServer() {
	listener := initTCPListener()
	go runAcceptLoop(listener)
	monitorForGracefulShutdown(listener)
}

func initTCPListener() net.Listener {
	serverConfig := config.GetServerConfig()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", serverConfig.Port))
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	logger.InfoLogger.Print("Server is listening on port " + serverConfig.Port)

	return listener
}

func runAcceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept() // TODO: set deadline?
		if err != nil {
			if isShutdownInitiated {
				break
			}

			logger.ErrorLogger.Print(err)

			continue
		}

		waitGroup.Add(1)

		go handleConn(conn)
	}
}

func monitorForGracefulShutdown(listener io.Closer) {
	signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quitChannel

	isShutdownInitiated = true

	listener.Close()
	waitGroup.Wait()
	services.TearDown()
	logger.InfoLogger.Println("Server successfully shutdown")
	os.Exit(0)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	defer waitGroup.Done()

	//for {
	//	_, err := io.WriteString(conn, time.Now().Format("15:04:05\n"))
	//	if err != nil {
	//		return // e.g., client disconnected
	//	}
	//	time.Sleep(1 * time.Second)
	//}
}
