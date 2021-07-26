.PHONY: server 
server:
	export TCP_ADDR=192.168.0.1:8080
	go run cmd/server/*.go


.PHONY: client
client:
	export TCP_ADDR=192.168.0.1:8080
	go run cmd/cleint/client.go