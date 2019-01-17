.PHONY: build

build:
	@go build -o bin/restgen -ldflags "-X main.version=$(git describe --tags --abbrev=0) -X main.revision=$(git rev-parse --short HEAD)" cmd/*
