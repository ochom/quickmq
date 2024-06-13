dev:
	@air

tidy:
	@go mod tidy
	
lint:
	@echo "Running linter..."
	@golangci-lint run

ui:
	@cd web && bun run build

docker:
	@docker build -t quickmq:latest .