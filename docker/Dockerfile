FROM golang:latest AS build

ARG SVC
EXPOSE 8080

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download -x

COPY . .

WORKDIR /build/cmd/${SVC} 

RUN go build

RUN cp ${SVC} ./entrypoint

ENTRYPOINT ./entrypoint
