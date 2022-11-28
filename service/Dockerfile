FROM golang:1.17 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .


RUN go build .


