.PHONY: all
all: 
	cd docker; docker-compose build && docker-compose up

.PHONY: all-mongo
all-mongo: 
	cd docker; docker-compose build && docker-compose --profile mongo up

.PHONY: all-redis
all-mongo: 
	cd docker; docker-compose build && docker-compose --profile redis up

.PHONY: all-postgres
all-mongo: 
	cd docker; docker-compose build && docker-compose --profile postgres up

.PHONY: tcpserver 
tcpserver:
	export TCP_ADDR=127.0.0.1:8080
	go run cmd/tcpserver/*.go

.PHONY: tcpclient
tcpclient:
	export TCP_ADDR=127.0.0.1:8080
	go run cmd/client/main.go

.PHONY: test 
test:
	go test -v ./pkg/httpserver/handler


