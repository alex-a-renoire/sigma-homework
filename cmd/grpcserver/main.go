package main

import (
	"log"
	"net"
	"os"

	"github.com/alex-a-renoire/sigma-homework/pkg/grpcserver"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/redisstorage"
	"google.golang.org/grpc"
)

type config struct {
	TCPport string
}

func getOsVars() *config {
	tcpPort := os.Getenv("GRPC_LISTEN_ADDRESS")
	if tcpPort == "" {
		tcpPort = ":50051"
	}

	return &config{
		TCPport:       tcpPort,
		DBType:        DBType,
		PGAddress:     pgAddress,
		RedisAddress:  redisAddress,
		RedisPassword: redisPassword,
		RedisDb:       db,
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

	//TODO: сделать слой контроллера - http или tcp - бизнес логика не должна меняться в зависимости от БД или GRPC
	log.Print(cfg.DBType)

	//create storage
	var db storage.Storage
	switch cfg.DBType {
	case "postgres":
		db, err = pgstorage.New(cfg.PGAddress)
		if err != nil {
			log.Printf("failed to connect to db: %s", err)
			return
		}
	case "redis":
		db = redisstorage.NewRDS(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDb)
	}

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterStorageServiceServer(s, grpcserver.New(db))

	log.Println("GRPC storage server starting...")
	//start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
