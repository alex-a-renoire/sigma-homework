package main

import (
	"log"
	"net"

	"github.com/alex-a-renoire/sigma-homework/pkg/grpcserver"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/mongostorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/redisstorage"
	"google.golang.org/grpc"
)

func main() {
	//Get configs
	cfg := getOsVars()

	//start listening on tcp
	lis, err := net.Listen("tcp", cfg.TCPport)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	log.Printf("DB type: %s", cfg.DBType)

	//create storage
	var db storage.Storage
	switch cfg.DBType {
	case "postgres":
		db, err = pgstorage.New(cfg.PGAddress)
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err)
		}
	case "redis":
		db = redisstorage.NewRDS(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDb)
	case "mongo":
		db, err = mongostorage.New(cfg.MongoAddress, cfg.MongoUser, cfg.MongoPassword)
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err)
		}
	}

	//create GRPC server
	s := grpc.NewServer()
	pb.RegisterStorageServiceServer(s, grpcserver.NewGRPC(db))

	log.Println("GRPC storage server starting...")
	//start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
