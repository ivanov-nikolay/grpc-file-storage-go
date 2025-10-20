GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=grpc-file-storage
BINARY_UNIX=$(BINARY_NAME)_unix

PROJECT_ROOT=github.com/grpc-file-storage-go
CMD_PATH=cmd/server
PROTO_PATH=api/proto

DOCKER_IMAGE=grpc-file-storage
DOCKER_TAG=latest

.PHONY: all build clean test run deps proto docker-build docker-run

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o bin/$(BINARY_NAME) ./$(CMD_PATH)

build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)_linux ./$(CMD_PATH)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf storage/

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

run:
	@echo "Running application..."
	$(GOBUILD) -o bin/$(BINARY_NAME) ./$(CMD_PATH)
	./bin/$(BINARY_NAME)

deps:
	@echo "Installing dependencies..."
	$(GOGET) -u google.golang.org/grpc
	$(GOGET) -u google.golang.org/protobuf
	$(GOGET) -u github.com/lib/pq
	$(GOGET) -u github.com/stretchr/testify
	$(GOGET) -u golang.org/x/sync

proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. $(PROTO_PATH)/*.proto

generate:
	@echo "Running go generate..."
	go generate ./...

lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	docker run -p 50051:50051 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up:
	@echo "Starting with Docker Compose..."
	docker-compose up --build

docker-compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  clean          - Clean build files"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run            - Build and run the application"
	@echo "  deps           - Install dependencies"
	@echo "  proto          - Generate protobuf code"
	@echo "  generate       - Run go generate"
	@echo "  lint           - Run linter"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start with Docker Compose"
	@echo "  docker-compose-down - Stop Docker Compose"
	@echo "  help           - Show this help message"