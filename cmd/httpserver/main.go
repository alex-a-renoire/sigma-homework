package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	grpccontroller "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/controller"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	httphandler "github.com/alex-a-renoire/sigma-homework/pkg/httpserver/handler"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/mongostorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/redisstorage"
	"github.com/alex-a-renoire/sigma-homework/service/csvservice"
	"github.com/alex-a-renoire/sigma-homework/service/personservice"
	"google.golang.org/grpc"
)

type config struct {
	HTTPAddr      string
	ConnType      string
	GRPCAddr      string
	DBType        string
	PGAddress     string
	RedisAddress  string
	RedisPassword string
	RedisDb       int
	MongoAddress  string
	MongoUser     string
	MongoPassword string
}

func getCfg() config {
	HTTPAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	if HTTPAddr == "" {
		HTTPAddr = ":8081"
	}

	ConnType := os.Getenv("CONN_TYPE")
	if ConnType == "" {
		ConnType = "grpc"
	}

	GRPCAddr := os.Getenv("HTTP_GRPC_ADDRESS")
	if GRPCAddr == "" {
		GRPCAddr = ":50051"
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

	mongoAddress := os.Getenv("MONGO_ADDRESS")
	if mongoAddress == "" {
		mongoAddress = ":27017"
	}

	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	if mongoUser == "" {
		mongoUser = "sigma-intern"
	}

	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoPassword == "" {
		mongoPassword = "sigma"
	}

	//possible values: postgres, redis, mongo
	DBType := os.Getenv("DB_TYPE")
	if DBType == "" {
		DBType = "mongo"
	}

	return config{
		HTTPAddr:      HTTPAddr,
		ConnType:      ConnType,
		GRPCAddr:      GRPCAddr,
		DBType:        DBType,
		PGAddress:     pgAddress,
		RedisAddress:  redisAddress,
		RedisPassword: redisPassword,
		RedisDb:       db,
		MongoAddress:  mongoAddress,
		MongoUser:     mongoUser,
		MongoPassword: mongoPassword,
	}
}

func main() {
	cfg := getCfg()

	log.Printf("DB type:" + cfg.DBType)

	//create storage
	var (
		storage personservice.PersonStorage
		err     error
	)

	if cfg.ConnType == "direct" {
		switch cfg.DBType {
		case "postgres":
			storage, err = pgstorage.New(cfg.PGAddress)
			if err != nil {
				log.Printf("failed to connect to db: %s", err)
				return
			}
		case "redis":
			storage = redisstorage.NewRDS(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDb)
		case "mongo":
			storage, err = mongostorage.New(cfg.MongoAddress, cfg.MongoUser, cfg.MongoPassword)
			if err != nil {
				log.Printf("failed to connect to db: %s", err)
				return
			}
		}
	} else if cfg.ConnType == "grpc" {
		//create storage service
		conn, err := grpc.Dial(cfg.GRPCAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect to grpc: %v", err)
		}
		defer conn.Close()

		storage = grpccontroller.New(pb.NewStorageServiceClient(conn))
	}

	//create service with storage
	personservice := personservice.New(storage)

	csvservice := csvservice.New(personservice)

	//create handler with controller
	sh := httphandler.New(personservice, *csvservice)

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
