package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"google.golang.org/grpc"

	httphandler "github.com/alex-a-renoire/sigma-homework/pkg/httpserver/handler"
	"github.com/alex-a-renoire/sigma-homework/service"
)

type config struct {
	HTTPAddr string
	GRPCAddr string
}

func getCfg() config {
	HTTPAddr := os.Getenv("HTTP_ADDR")
	if HTTPAddr == "" {
		HTTPAddr = ":8081"
	}

	GRPCAddr := os.Getenv("GRPC_ADDR")
	if GRPCAddr == "" {
		GRPCAddr = ":50051"
	}

	return config{
		HTTPAddr: HTTPAddr,
		GRPCAddr: GRPCAddr,
	}
}

func main() {
	cfg := getCfg()

	//create storage service
	conn, err := grpc.Dial(cfg.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	storageClient := pb.NewStorageServiceClient(conn)

	//create service with storage
	service := service.NewGRPC(storageClient)

	//create handler with controller
	sh := httphandler.New(service)

	srv := http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: sh.GetRouter(),
	}

	//graceful shutdown of server
	sigC := make(chan os.Signal, 1)
	defer close(sigC)
	go func() {
		<-sigC
		srv.Shutdown(context.TODO())
	}()
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	//Start the server
	log.Print("Starting server...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("error: http server failed: %s", err)
	}
}
