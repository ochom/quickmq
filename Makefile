dev:
	@air

tidy:
	@go mod tidy
	
lint:
	@echo "Running linter..."
	@golangci-lint run
