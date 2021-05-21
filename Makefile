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

.PHONY: setup
setup: ## Run setup scripts to prepare development environment
	chmod 777 -R scripts/
	@scripts/setup.sh

.PHONY: lint
lint:
	golangci-lint run