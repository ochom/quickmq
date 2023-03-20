SHELL=/bin/bash

dev:
	air

build:
	@echo "building dev ..."
	@go build -race  -o dist/pub examples/publisher/main.go
	@go build -race  -o dist/sub examples/consumer/main.go

pub:
	@echo "Running publisher ..."
	@./dist/pub

sub:
	@echo "Running consumer ..."
	@./dist/sub

lint:
	@echo "Linting..."
	@golangci-lint run