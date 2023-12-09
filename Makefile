BIN = $(shell go env GOPATH)/bin
run:
	go run cmd/main.go

lint:
	$(BIN)/golint ./...

build:
	go build -o bin/cube-dashboard-api cmd/main.go

all: lint build