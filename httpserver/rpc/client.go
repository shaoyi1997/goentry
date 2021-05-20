package rpc

import (
	"net"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/config"
	"google.golang.org/protobuf/proto"
)

type IRPCClient interface {
	CallMethod(pb.RpcRequest_Method, interface{}, interface{}) error
	Close()
}

type Client struct {
	connectionPool IPool
}

func NewRPCClient() (IRPCClient, error) {
	serverConfig := config.GetServerConfig()
	poolConfig := config.GetPoolConfig()

	pool, err := NewPool(&PoolConfig{
		InitCap:     poolConfig.InitialCapacity,
		MaxCap:      poolConfig.MaxCapacity,
		WaitTimeout: time.Duration(poolConfig.InitialCapacity) * time.Second,
		Factory:     func() (net.Conn, error) { return net.Dial("tcp", serverConfig.TCPAddress) },
	})
	if err != nil {
		return nil, err
	}

	return &Client{connectionPool: pool}, nil
}

func (rpcClient *Client) Close() {
	rpcClient.connectionPool.Close()
}

func (rpcClient *Client) CallMethod(method pb.RpcRequest_Method, requestMessage interface{},
	response interface{}) error {
	serializedRequest, err := rpc.SerializeMessage(method, requestMessage)
	if err != nil {
		return err
	}

	connection, err := rpcClient.connectionPool.Get()
	if err != nil {
		return err
	}
	defer connection.Close()

	if err = sendRequest(connection, serializedRequest); err != nil {
		return err
	}

	return receiveResponse(connection, response)
}

func receiveResponse(conn net.Conn, response interface{}) error {
	messageBuffer, err := rpc.ReadMessageBufferFromConnection(conn)
	if err != nil {
		logger.ErrorLogger.Println("Failed to read response body", err)

		return err
	}

	// Reads from the 5th byte onwards. Ignore the method parameter
	if err := proto.Unmarshal(messageBuffer[4:], response.(proto.Message)); err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal response", err)

		return err
	}

	return nil
}

func sendRequest(connection net.Conn, message []byte) error { //nolint:interfacer
	if _, err := connection.Write(message); err != nil {
		logger.WarningLogger.Println("Failed to send RPC request", err)

		return err
	}

	return nil
}
