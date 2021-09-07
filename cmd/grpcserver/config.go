package main

import (
	"os"
	"strconv"
)

type config struct {
	TCPport        string
	DBType         string
	PGAddress      string
	RedisAddress   string
	RedisPassword  string
	RedisDb        int
	MongoAddress   string
	MongoUser      string
	MongoPassword  string
	ElasticAddress string
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

	elasticAddress := os.Getenv("ELASTIC_ADDRESS")
	if elasticAddress == "" {
		elasticAddress = ":9200"
	}

	//possible values: postgres, redis, mongo, elastic
	DBType := os.Getenv("DB_TYPE")
	if DBType == "" {
		DBType = "postgres"
	}

	return &config{
		TCPport:        tcpPort,
		DBType:         DBType,
		PGAddress:      pgAddress,
		RedisAddress:   redisAddress,
		RedisPassword:  redisPassword,
		RedisDb:        db,
		MongoAddress:   mongoAddress,
		MongoUser:      mongoUser,
		MongoPassword:  mongoPassword,
		ElasticAddress: elasticAddress,
	}
}
