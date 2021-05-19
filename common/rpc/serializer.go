package rpc

import (
	"bytes"
	"encoding/binary"
	"net"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"google.golang.org/protobuf/proto"
)

func SerializeMessage(method pb.RpcRequest_Method, args interface{}) ([]byte, error) {

	message, err := prepareRequestMessage(method, args)
	if err != nil {
		return nil, err
	}

	fullMessage, err := prependLenToMessage(message)
	if err != nil {
		return nil, err
	}

	return fullMessage.Bytes(), nil
}

func prepareRequestMessage(method pb.RpcRequest_Method, args interface{}) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(nil)

	encodedArgs, err := proto.Marshal(args.(proto.Message))
	if err != nil {
		logger.ErrorLogger.Println("Failed to encode request args:", err)
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, method)
	if err != nil {
		logger.ErrorLogger.Println("Failed to write request method: ", err)
		return nil, err
	}

	err = binary.Write(buffer, nil, encodedArgs)
	if err != nil {
		logger.ErrorLogger.Println("Failed to write request args:", err)
		return nil, err
	}
	return buffer, nil
}

func prependLenToMessage(message *bytes.Buffer) (*bytes.Buffer, error) {
	length := uint32(message.Len())
	fullMessage := bytes.NewBuffer(nil)

	err := binary.Write(fullMessage, binary.BigEndian, length)
	if err != nil {
		logger.ErrorLogger.Println("Failed to write length: ", err)
		return nil, err
	}

	err = binary.Write(fullMessage, nil, message.Bytes())
	if err != nil {
		logger.ErrorLogger.Println("Failed to write message: ", err)
		return nil, err
	}
	return fullMessage, err
}

func ReadMessageBufferFromConnection(conn net.Conn) ([]byte, error) {
	var lenByte [4]byte
	if _, err := conn.Read(lenByte[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenByte[:])

	buffer := make([]byte, length)
	var totalLengthRead uint32

	for {
		read, err := conn.Read(buffer[totalLengthRead:])
		totalLengthRead += uint32(read)
		if totalLengthRead >= length {
			return buffer, nil
		} else if err != nil {
			return nil, err
		}
	}
}
