package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"google.golang.org/protobuf/proto"

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
	logger.InitLogger()
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
	deserializeRequest(conn)
}

func deserializeRequest(conn net.Conn) {
	for {
		messageBuffer, err := rpc.ReadMessageBufferFromConnection(conn)
		if err != nil {
			if err == io.EOF {
				logger.ErrorLogger.Println("Failed to read length of request", err)
				return
			}
			continue
		}

		var responseBuffer bytes.Buffer
		method := binary.BigEndian.Uint32(messageBuffer[:4])

		switch method {
		case uint32(pb.RpcRequest_Login):
			logger.InfoLogger.Println("Login request")
			var args pb.LoginRequest
			err := proto.Unmarshal(messageBuffer[4:], &args)
			if err != nil {
				logger.ErrorLogger.Println(err)
				return
			}
			var userId uint64 = 2
			response := &pb.LoginResponse{
				UserId:   &userId,
				Username: args.Username,
				Password: args.Password,
				Token:    args.Username,
			}
			message, err := rpc.SerializeMessage(pb.RpcRequest_Login, response)
			if err != nil {
				logger.ErrorLogger.Println("Failed to serialize message:", err)
				return
			}

			responseBuffer.Write(message)
		}
		conn.Write(responseBuffer.Bytes())
	}

}
