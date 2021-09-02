# A TCP client-server app for Persons | HTTP server for Persons

It lets clients write data to the database. The basic entity is a person with a certain name. 
Clients accept data in JSON format, e.g.:

- AddPerson {"func_name":"AddPerson", "data":{"name":"Bob"}}
- GetPerson {"func_name":"GetPerson", "data":{"id":1}}
- GetAllPersons {"func_name":"GetAllPersons", "data":{}}
- UpdatePerson {"func_name":"UpdatePerson", "data":{"id":1, "name":"Alice"}}
- DeletePerson {"func_name":"DeletePerson", "data":{"id":1}}

## Cases accounted for: 

- Input is not a JSON
- Wrong data type of a json field: {"func_name":"GetPerson", "data":{"id":"1"}}
- Wrong json field tag name / absence of a required field: {"func":"GetPerson", "data":{"id":1}} / {"func_name":"GetPerson"} / {"data":{"id":0}}
- Wrong field value {"func_name":"wrong_func", "data":{"name":"Bob"}}
- Delete / get a person with non-existent ID

## How to run tcp-app

make tcpserver
make tcpclient
make test

## How to run http-app

make httpserver

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