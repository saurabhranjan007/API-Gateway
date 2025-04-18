# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o zeneye-gateway ./cmd/zeneye-gateway

EXPOSE 8080

CMD ["./zeneye-gateway"]
