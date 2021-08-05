package main

import (
	"log"
	"net"
	"os"

	"github.com/alex-a-renoire/sigma-homework/pkg/grpcserver"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/inmemory"
	"google.golang.org/grpc"
)

type config struct {
	TCPport string
}

func getOsVars() *config {
	tcpPort := os.Getenv("POSTAL_LISTEN_ADDRESS")
	if tcpPort == "" {
		tcpPort = ":50051"
	}

	return &config{
		TCPport: tcpPort,
	}
}

func main() {
	//Get configs for it
	cfg := getOsVars()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TCPport)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	//create storage
	db := inmemory.New()

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterStorageServiceServer(s, &grpcserver.StorageServer{
		DB: db,
	})

	//start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
