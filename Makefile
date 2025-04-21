.PHONY: build run test clean swag

# Application name
APP_NAME=healthcare-app

# Build the application
build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

# Run the application
run:
	go run ./cmd/api/main.go

# Run tests
test:
	go test -v ./internal/...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

# Generate Swagger documentation
swag:
	swag init -g cmd/api/main.go -o docs

# Clean build artifacts
clean:
	rm -rf ./bin

# Create database migration
migrate-up:
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/healthcare?sslmode=disable" up

# Rollback database migration
migrate-down:
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/healthcare?sslmode=disable" down

# Create a new migration file
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Setup initial development environment
setup:
	@echo "Creating necessary directories..."
	@mkdir -p bin migrations
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Setup complete!"

# Default command
default: run 