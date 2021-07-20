package main

import (
	"os"
	"os/signal"
	"syscall"

	dummytcp "github.com/alex-a-renoire/tcp"
	"github.com/alex-a-renoire/tcp/storage/inmemory"
	"github.com/alex-a-renoire/tcp/tcpserver/handler"
	"github.com/alex-a-renoire/tcp/tcpserver/server"
)

func main() {
	//create a storage
	s := inmemory.New()

	//create a request handler
	h := handler.New(&s)

	//create a server
	srv := server.New(dummytcp.TCP_ADDR, h)

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
