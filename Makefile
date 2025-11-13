.PHONY: build build-all clean run test

APP_NAME=miner
BUILD_DIR=build
CMD_DIR=cmd/miner

build:
	@echo "Building for current platform..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)

build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for macOS (AMD64)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./$(CMD_DIR)
	
	@echo "Building for macOS (ARM64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./$(CMD_DIR)
	
	@echo "Building for Linux (AMD64)..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./$(CMD_DIR)
	
	@echo "Building for Windows (AMD64)..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./$(CMD_DIR)
	
	@echo "Build complete! Binaries in $(BUILD_DIR)/"

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	rm -rf bin

run:
	@echo "Running $(APP_NAME)..."
	go run ./$(CMD_DIR)

test:
	@echo "Running tests..."
	go test -v ./...

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run || true
