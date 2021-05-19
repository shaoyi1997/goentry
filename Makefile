MAKEFLAGS += -j3

run: tcpserver fileserver httpserver

.PHONY: tcpserver
tcpserver:
	go run tcpserver/tcpserver.go

.PHONY: fileserver
fileserver:
	go run fileserver/fileserver.go

.PHONY: httpserver
httpserver:
	sleep 3 # to allow time for tcpserver to start up
	go run httpserver/httpserver.go
