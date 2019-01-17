.PHONY: build test

build:
	@go build -v -o bin/restgen -ldflags "-X main.version=$(shell git describe --tags --abbrev=0) -X main.revision=$(shell git rev-parse --short HEAD)" cmd/*

test:
	@go test -v ./...
