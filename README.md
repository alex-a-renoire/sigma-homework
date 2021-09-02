# Go RESTful API dummy project

This package serves the goal of studying golang step by step and is implemented by myself in oder to learn various design patterns and libraries. 

The basic entity the project deals with is person. For the moment, the following features are supported:

- RESTful endpoints;
- CRUD operations (mongodb, postgers, redis);
- Upload and download entities in the CSV format;
- JWT-based pseudo-authentication;
- Environment dependent application configuration management;
- Structured logging with contextual information;
- Error handling with proper error response generation;
- Data validation;

The project uses the following go packages: 

- Redis: <a href="github.com/go-redis/redis"go-redis/redis</a
- JWT: <a href="github.com/golang-jwt/jwt"golang-jwt/jwt</a
- UUID: <a href="github.com/google/uuid"google/uuid</a
- Routing: <a href="github.com/gorilla/mux"gorilla/mux</a
- CSV handling: <a href="github.com/jszwec/csvutil"jszwez/csvutil</a
- Postgres: <a href="github.com/lib/pq"lib/pq</a
- Mongo: <a href="go.mongodb.org/mongo-driver"mongo-driver</a
- GRPC: <a href="google.golang.org/grpc"grpc</a

# Getting started

First, make sure you have *docker* and *docker-compose* installed on your system. To run the product with default configuration, open the root folder in bash and type:

 make all

To run the project with a cartain database (mongo, redis, postgres), type:

 make all-mongo 

or

 make all-postgres

etc.

**Important!** You have to also specify the database type in the *docker/.env* file (DB_TYPE). 

You can also select the type of storage connection (*remote* or *grpc*) (CONN_TYPE)

By default, RESTful API server runs in a container at :8081. GRPC runs at :50051, redis at :6379, postgres at :5432 and mongo at :27017. The API provides the following endpoints:

- `GET /persons` : Requests all persons from the database

- `POST /persons` : Posts a person to the database. Returns its id.

- `GET /persons` : Requests a person with a specified id

- `PUT /persons/{id}` : Updates a person with a specified id

- `DELETE /persons/{id}` : Deletes the person with a specified id

- `GET /persons/me` : Looks at the JWT token and tells the user who she is

- `GET /persons/upload` : Renders a webpage with a CSV document upload form

- `POST /persons/upload` : Uploads the CSV document and saves the persons from the document to the database

- `GET /login/{id}` : Generates a JWT token for the session

The response format is JSON or a byte array in case of a JWT token.

If you have <a href="https://httpie.io/"httpie</a or some other API client tool (e.g. Postman), you may try the following more complex scenarios:

```
# post a user
http -v POST 127.0.0.1:8081/persons "Name"="Bob"
# should return a JSON like {id: 7c7650fe-843c-476e-8132-ce754e15314c, name: Alice}

# get the list of users
http -v GET 127.0.0.1:8081/persons
# should return an array of JSONs like [{id: 7c7650fe-843c-476e-8132-ce754e15314c, name: Alice}]

# get some certain user
http -v GET 127.0.0.1:8081/persons/7c7650fe-843c-476e-8132-ce754e15314c
# should return a JSON like {id: 7c7650fe-843c-476e-8132-ce754e15314c, name: Alice}

# delete some certain user
http -v DELETE 127.0.0.1:8081/persons/7c7650fe-843c-476e-8132-ce754e15314c
# should return 200 OK

# update some certain user
http -v PUT 127.0.0.1:8081/persons/7c7650fe-843c-476e-8132-ce754e15314c "Name"="Josh"
# should return 200 OK

# Get current session user user from JWT token 
http -v GET 127.0.0.1:8081/persons/me 'Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzA1ODI1NzIsImlhdCI6MTYzMDU4MDc3MiwiSWQiOiI3Yzc2NTBmZS04NDNjLTQ3NmUtODEzMi1jZTc1NGUxNTMxNGMiLCJlbWFpbCI6IkJvYiJ9.4dr4kNWuKUiVIFxAv8v_fBmgWUOVopmnw7-NTApRWIU'
# should return a JSON like {id: 7c7650fe-843c-476e-8132-ce754e15314c, name: Alice}
```


Try the URLs http://localhost:8081/persons/download or http://localhost:8081/persons/upload in a browser to download or upload CSV files.

# Project layout
```
├── cmd                         main applications of the project
│   ├── grpcserver              grpc server which features database functions
│   │   ├── config.go           server configuration (from env file)
│   │   └── main.go             
│   ├── httpserver              http server 
│   │   ├── config.go           
│   │   └── main.go
│   └── mem                     old tcp version of the project
│       ├── client
│       │   └── main.go
│       └── tcpserver
│           └── main.go
├── docker                      
│   ├── docker-compose.yaml     
│   ├── Dockerfile              all the services are dockerized
│   └── templates
│       └── upload.html         template for the csv upload form
├── go.mod
├── go.sum
├── Makefile
├── model
│   ├── action.go               is used only in old tcp version
│   ├── errors.go
│   ├── personaddupdate.go      models for different logic of the person
│   └── person.go
├── pkg
│   ├── grpcserver
│   │   ├── controller          
│   │   │   └── controller.go   
│   │   ├── proto
│   │   │   ├── service_grpc.pb.go
│   │   │   ├── service.pb.go
│   │   │   └── service.proto
│   │   └── server.go
│   ├── httpserver
│   │   └── handler
│   │       ├── httphandler.go
│   │       └── httphandler_test.go
│   ├── storage
│   │   ├── inmemory
│   │   │   └── inmemory.go
│   │   ├── mockstorage.go
│   │   ├── mongostorage
│   │   │   ├── init-mongo.js
│   │   │   └── mongostorage.go
│   │   ├── pgstorage
│   │   │   ├── pgstorage.go
│   │   │   └── schema.sql
│   │   ├── redisstorage
│   │   │   └── redisstorage.go
│   │   └── storage.go
│   └── tcpserver
│       ├── controller
│       │   └── controller.go
│       ├── handler
│       │   └── handler.go
│       └── server.go
├── README.md
├── roadmap.txt
└── service
    ├── authservice
    │   └── service.go
    ├── csvservice
    │   ├── service.go
    │   ├── service_test.go
    │   └── testdata
    │       ├── add.csv
    │       ├── empty.csv
    │       ├── emptyfields.csv
    │       ├── onlyheaders.csv
    │       ├── rename.csv
    │       ├── toomanyfields.csv
    │       └── wrongid.csv
    └── personservice
        ├── service.go
        └── service_test.go
```
