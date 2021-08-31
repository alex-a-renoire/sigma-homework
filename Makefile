.PHONY: all
all: 
	cd docker; docker-compose build && docker-compose up

.PHONY: all-mongo
all-mongo: 
	cd docker; docker-compose build && docker-compose --profile mongo up

.PHONY: tcpserver 
tcpserver:
	export TCP_ADDR=127.0.0.1:8080
	go run cmd/tcpserver/*.go


.PHONY: tcpclient
tcpclient:
	export TCP_ADDR=127.0.0.1:8080
	go run cmd/client/main.go

.PHONY: httpserver 
httpserver:
	export HTTP_ADDR=127.0.0.1:8081
	go run cmd/httpserver/*.go

.PHONY: test 
test:
	go test -v ./pkg/httpserver/handler


