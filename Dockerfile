# syntax=docker/dockerfile:1

FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
EXPOSE 8080