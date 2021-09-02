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

- Redis: <a href="github.com/go-redis/redis">go-redis/redis</a>
- JWT: <a href="github.com/golang-jwt/jwt">golang-jwt/jwt</a>
- UUID: <a href="github.com/google/uuid">google/uuid</a>
- Routing: <a href="github.com/gorilla/mux">gorilla/mux</a>
- CSV handling: <a href="github.com/jszwec/csvutil">jszwez/csvutil</a>
- Postgres: <a href="github.com/lib/pq">lib/pq</a>
- Mongo: <a href="go.mongodb.org/mongo-driver">mongo-driver</a>
- GRPC: <a href="google.golang.org/grpc">grpc</a>

## Getting started

First, make sure you have *docker* and *docker-compose* installed on your system. To run the product with default configuration, open the root folder in bash and type:

> make all

To run the project with a cartain database (mongo, redis, postgres), type:

> make all-mongo 

or

> make all-postgres

etc.

**Important!** You have to also specify the database type in the *docker/.env* file (DB_TYPE). 

You can also select the type of storage connection (*remote* or *grpc*) (CONN_TYPE)

By default, RESTful API server runs in a container at :8081. GRPC runs at :50051, redis at :6379, postgres at :5432 and mongo at :27017. The API provides the following endpoints:

- <mark>GET /persons :</mark> Requests all persons from the database

- <mark>POST /persons :</mark> Posts a person to the database. Returns its id.

- <mark>GET /persons :</mark> Requests a person with a specified id

- <mark>PUT /persons/{id} :</mark> Updates a person with a specified id

- <mark>DELETE /persons/{id} :</mark> Deletes the person with a specified id

- <mark>GET /persons/me :</mark> Looks at the JWT token and tells the user who she is

- <mark>GET /persons/upload :</mark> Renders a webpage with a CSV document upload form

- <mark>POST /persons/upload :</mark> Uploads the CSV document and saves the persons from the document to the database

- <mark>GET /login/{id} :</mark> Generates a JWT token for the session

The response format is JSON or a byte array in case of a JWT token.

If you have <a href="https://httpie.io/">httpie</a> or some other API client tool (e.g. Postman), you may try the following more complex scenarios:

> #post a user
> http -v POST 127.0.0.1:8081/persons "Name"="Bob"
> #should return a JSON like {id: sdfsdfsdf, name: Alice}
>
>

Try the URL http://localhost:8081/persons in a browser, and you should see something like "OK v1.0.0" displayed.

## Cases accounted for: 

- Input is not a JSON
- Wrong data type of a json field: {"func_name":"GetPerson", "data":{"id":"1"}}
- Wrong json field tag name / absence of a required field: {"func":"GetPerson", "data":{"id":1}} / {"func_name":"GetPerson"} / {"data":{"id":0}}
- Wrong field value {"func_name":"wrong_func", "data":{"name":"Bob"}}
- Delete / get a person with non-existent ID


## how to test http-app


- **REQUEST** http -v POST 127.0.0.1:8081/persons "Name"="Bob"      **RESPONSE**: Person with id 1 and name Bob added
- **REQUEST** http -v POST 127.0.0.1:8081/persons "Name"="Alice"    **RESPONSE**: Person with id 2 and name Alice added
- **REQUEST** http -v GET 127.0.0.1:8081/persons                    **RESPONSE**: All persons in the storage are [{1 Bob} {2 Alice}]
- **REQUEST** http -v GET 127.0.0.1:8081/persons/1                  **RESPONSE**: Person with id 1 has name Bob
- **REQUEST** http -v GET 127.0.0.1:8081/persons/3                  **RESPONSE**: there is no such record   
- **REQUEST** http -v DELETE 127.0.0.1:8081/persons/1               **RESPONSE**: Person with id 1 deleted
- **REQUEST** http -v GET 127.0.0.1:8081/persons/1                  **RESPONSE**: error: person with id 1 not found
- **REQUEST** http -v GET 127.0.0.1:8081/persons                    **RESPONSE**: All persons in the storage are [{2 Alice}]
- **REQUEST** http -v PUT 127.0.0.1:8081/persons/2 "Name"="Rachel" **RESPONSE**: Person with id 2 updated with name Rachel
- **REQUEST** http -v GET 127.0.0.1:8081/persons                    **RESPONSE**: All persons in the storage are [{2 Rachel}]

- **REQUEST** http -v GET 127.0.0.1:8081/persons/dump
- **REQUEST** 127.0.0.1:8081/persons/upload
- **REQUEST** http -v GET 127.0.0.1:8081/login/7c7650fe-843c-476e-8132-ce754e15314c

http -v GET 127.0.0.1:8081/persons/myuser 'Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzA1ODI1NzIsImlhdCI6MTYzMDU4MDc3MiwiSWQiOiI3Yzc2NTBmZS04NDNjLTQ3NmUtODEzMi1jZTc1NGUxNTMxNGMiLCJlbWFpbCI6IkJvYiJ9.4dr4kNWuKUiVIFxAv8v_fBmgWUOVopmnw7-NTApRWIU'