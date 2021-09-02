package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-a-renoire/sigma-homework/pkg/storage/inmemory"
	"github.com/alex-a-renoire/sigma-homework/pkg/tcpserver"
	tcpcontroller "github.com/alex-a-renoire/sigma-homework/pkg/tcpserver/controller"
	"github.com/alex-a-renoire/sigma-homework/pkg/tcpserver/handler"
	"github.com/alex-a-renoire/sigma-homework/service/authservice"
	personservice "github.com/alex-a-renoire/sigma-homework/service/personservice"
)

type config struct {
	TCPAddr string
}

func getCfg() config {
	TCPAddr := os.Getenv("TCP_ADDR")
	if TCPAddr == "" {
		TCPAddr = "127.0.0.1:8080"
	}

	return config{
		TCPAddr: TCPAddr,
	}
}

func main() {
	cfg := getCfg()

	//create a storage
	s := inmemory.New()

	//create a service
	authService := authservice.New("some_unsafe_secret")

	service := personservice.New(s, authService)

	//create a controller for transport layer
	controller := tcpcontroller.New(service)

	//create a request handler
	h := handler.New(controller)

	//create a server
	srv := tcpserver.New(cfg.TCPAddr, h)

	//Graceful shutdown
	sigC := make(chan os.Signal, 1)
	defer close(sigC)

	go func() {
		<-sigC
		srv.Stop()
	}()

	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	//run the server
	srv.Wg.Add(1)
	go srv.Serve()

	srv.Wg.Wait()
}
