package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services"
)

var (
	quitChannel         = make(chan os.Signal)
	waitGroup           = sync.WaitGroup{}
	isShutdownInitiated = false
	service             *services.Services
)

func main() {
	logger.InitLogger()
	config.InitConfig()

	service = services.Init()

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

	logger.InfoLogger.Println("TCP server is listening on port:", serverConfig.Port)

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

	for {
		messageBuffer, err := rpc.ReadMessageBufferFromConnection(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			continue
		}

		responseMessage, err := routeRequest(messageBuffer)
		if err != nil {
			logger.ErrorLogger.Println("Failed to handle request:", err)

			return
		}

		_, err = conn.Write(responseMessage)
		if err != nil {
			logger.ErrorLogger.Println("Failed to write response:", err)

			return
		}
	}
}

func routeRequest(messageBuffer []byte) ([]byte, error) {
	var (
		responseMessage []byte
		err             error
	)

	method := binary.BigEndian.Uint32(messageBuffer[:4])
	message := messageBuffer[4:]

	switch method {
	case uint32(pb.RpcRequest_Login):
		responseMessage, err = service.User.Login(message)
	case uint32(pb.RpcRequest_Update):
		responseMessage, err = service.User.Update(message)
	case uint32(pb.RpcRequest_Register):
		responseMessage, err = service.User.Register(message)
	case uint32(pb.RpcRequest_GetUser):
		responseMessage, err = service.User.GetUser(message)
	case uint32(pb.RpcRequest_Logout):
		responseMessage, err = service.User.Logout(message)
	}

	if err != nil {
		logger.ErrorLogger.Println("Failed to execute ", method, " request:", err)

		return nil, err
	}

	return responseMessage, err
}
