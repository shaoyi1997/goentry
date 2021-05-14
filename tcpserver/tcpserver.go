package main

import (
	"fmt"
	"net"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services"
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	services.InitDB()
	initTCPServer()
}

func initTCPServer() {
	serverConfig := config.GetServerConfig()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", serverConfig.Port))
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	logger.InfoLogger.Print("Server is listening on port " + serverConfig.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.ErrorLogger.Print(err)

			continue
		}

		go handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
}
