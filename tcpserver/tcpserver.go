package main

import (
	"fmt"
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	config2 "git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	_ "git.garena.com/shaoyihong/go-entry-task/tcpserver/services"
	"net"
)

func main() {
	config := config2.GetServerConfig()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
	logger.InfoLogger.Print("Server is listening on port " + config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.ErrorLogger.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
}
