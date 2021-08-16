package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/alex-a-renoire/sigma-homework/pkg/grpcserver"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/redisstorage"
	"google.golang.org/grpc"
)

type config struct {
	TCPport       string
	DBType        string
	PGAddress     string
	RedisAddress  string
	RedisPassword string
	RedisDb       int
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

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "127.0.0.1:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	var (
		db  int
		err error
	)
	redisDb := os.Getenv("REDIS_DB")
	if redisDb != "" {
		db, err = strconv.Atoi(redisDb)
		if err != nil {
			panic(err)
		}
	}

	//possible values: postgres, redis, mongo
	DBType := os.Getenv("DB_TYPE")
	if DBType == "" {
		DBType = "postgres"
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
		log.Print("redis")
		log.Print(db)
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
