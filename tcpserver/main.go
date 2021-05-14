package main

import (
	config2 "git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	"log"
	"net"
)

func main() {
	config := config2.GetServerConfig()
	listener, err := net.Listen("tcp", ":" + config.Port)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Server is listening on port " + config.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

}
