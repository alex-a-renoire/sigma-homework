package main

import (
	"os"
	"strconv"
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
	JWTSecret     string
}

func getCfg() config {
	HTTPAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	if HTTPAddr == "" {
		HTTPAddr = ":8081"
	}

	//Initialize JWT sectret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "not_really_secure_secret"
	}

	//Connection type: direct or grpc
	ConnType := os.Getenv("CONN_TYPE")
	if ConnType == "" {
		ConnType = "grpc"
	}

	GRPCAddr := os.Getenv("HTTP_GRPC_ADDRESS")
	if GRPCAddr == "" {
		GRPCAddr = ":50051"
	}

	//Databases initialization: POSTGRES
	pgAddress := os.Getenv("PG_ADDRESS")
	if pgAddress == "" {
		pgAddress = "host=db port=5432 dbname=persons user=persons password=pass sslmode=disable"
	}

	//Databases initialization: REDIS
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

	//Databases initialization: MONGO
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
		JWTSecret:     jwtSecret,
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
