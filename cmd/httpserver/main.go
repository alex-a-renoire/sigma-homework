package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpccontroller "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/controller"
	pb "github.com/alex-a-renoire/sigma-homework/pkg/grpcserver/proto"
	httphandler "github.com/alex-a-renoire/sigma-homework/pkg/httpserver/handler"
	elastic "github.com/alex-a-renoire/sigma-homework/pkg/storage/elasticsearch"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/mongostorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/pgstorage"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage/redisstorage"
	"github.com/alex-a-renoire/sigma-homework/service/authservice"
	"github.com/alex-a-renoire/sigma-homework/service/csvservice"
	"github.com/alex-a-renoire/sigma-homework/service/personservice"
	"google.golang.org/grpc"
)

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
				log.Fatalf("failed to connect to db: %s", err)
			}
		case "redis":
			storage = redisstorage.NewRDS(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDb)
		case "mongo":
			storage, err = mongostorage.New(cfg.MongoAddress, cfg.MongoUser, cfg.MongoPassword)
			if err != nil {
				log.Fatalf("failed to connect to db: %s", err)
			}
		case "elastic":
			storage, err = elastic.New(cfg.ElasticAddress)
			if err != nil {
				log.Fatalf("failed to connect to db: %s", err)
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

	//create services with storage
	personservice := personservice.New(storage)

	authservice := authservice.New(personservice, cfg.JWTSecret)

	csvservice := csvservice.New(personservice)

	//create handler
	sh := httphandler.New(personservice, *csvservice, authservice)

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
