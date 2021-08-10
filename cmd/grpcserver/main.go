package main

import (
	"log"
	"net"
	"os"

	"github.com/alex-a-renoire/sigma-homework/pkg/grpcserver"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"google.golang.org/grpc"
)

type config struct {
	TCPport   string
	PGAddress string
}

func getOsVars() *config {
	tcpPort := os.Getenv("GRPC_LISTEN_ADDRESS")
	if tcpPort == "" {
		tcpPort = ":50051"
	}

	pgAddress := os.Getenv("PG_ADDRESS")
	if pgAddress == "" {
		pgAddress = "host=db port=5432 dbname=persons user=persons password=pass sslmode=disable"
	}

	return &config{
		TCPport:   tcpPort,
		PGAddress: pgAddress,
	}
}

func main() {
	//Get configs
	cfg := getOsVars()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TCPport)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	//create storage
	db, err := pgstorage.New(cfg.PGAddress)
	if err != nil {
		log.Printf("failed to connect to db: %s", err)
		return
	}

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterStorageServiceServer(s, &grpcserver.StorageServer{
		DB: db,
	})

	log.Println("GRPC storage server starting...")
	//start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
