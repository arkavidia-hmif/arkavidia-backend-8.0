# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY *.go ./

RUN go build -o /arkavidia-backend

EXPOSE 8080