.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/app main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t quant-trading-app .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker Compose
docker-stop:
	docker-compose down

# Run database migrations (if using a migration tool)
migrate:
	@echo "Run database migrations here"

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run