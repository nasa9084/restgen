.PHONY: hook build test

hook:
	@cp githooks/* .git/hooks/
	@chmod +x .git/hooks/*

build:
	@go-assets-builder assets -p assets -s /assets -o internal/pkg/assets/assets.go
	@go build -v -o bin/restgen -ldflags "-X main.version=$(shell git describe --tags --abbrev=0) -X main.revision=$(shell git rev-parse --short HEAD)" cmd/*

test:
	@go test -v ./...
