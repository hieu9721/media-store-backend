# Media Store Backend - Makefile
# Cross-platform build automation

APP_NAME := media-store-backend
VERSION := 1.0.0
BUILD_DIR := build
LDFLAGS := -s -w -X main.Version=$(VERSION)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Colors
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: all build clean test run help deps windows linux mac arm docker

# Default target
all: clean build-all

# Help command
help:
	@echo "$(BLUE)Media Store Backend - Build System$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@echo "  $(GREEN)make build$(NC)          - Build for current platform"
	@echo "  $(GREEN)make build-all$(NC)      - Build for all platforms"
	@echo "  $(GREEN)make windows$(NC)        - Build for Windows (amd64, 386, arm64)"
	@echo "  $(GREEN)make linux$(NC)          - Build for Linux (amd64, 386, arm, arm64)"
	@echo "  $(GREEN)make mac$(NC)            - Build for macOS (amd64, arm64)"
	@echo "  $(GREEN)make arm$(NC)            - Build for ARM platforms"
	@echo "  $(GREEN)make run$(NC)            - Run the application"
	@echo "  $(GREEN)make test$(NC)           - Run tests"
	@echo "  $(GREEN)make clean$(NC)          - Clean build artifacts"
	@echo "  $(GREEN)make deps$(NC)           - Download dependencies"
	@echo "  $(GREEN)make docker$(NC)         - Build Docker image"
	@echo "  $(GREEN)make release$(NC)        - Create release archives"
	@echo ""

# Build for current platform
build:
	@echo "$(BLUE)Building for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

# Build for all platforms
build-all: clean windows linux mac
	@echo "$(GREEN)All builds complete!$(NC)"

# Windows builds
windows:
	@echo "$(BLUE)=== Windows Builds ===$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✓ Windows amd64$(NC)"
	GOOS=windows GOARCH=386 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-386.exe .
	@echo "$(GREEN)✓ Windows 386$(NC)"
	GOOS=windows GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe .
	@echo "$(GREEN)✓ Windows arm64$(NC)"

# Linux builds
linux:
	@echo "$(BLUE)=== Linux Builds ===$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	@echo "$(GREEN)✓ Linux amd64$(NC)"
	GOOS=linux GOARCH=386 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-386 .
	@echo "$(GREEN)✓ Linux 386$(NC)"
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	@echo "$(GREEN)✓ Linux arm64$(NC)"
	GOOS=linux GOARCH=arm $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm .
	@echo "$(GREEN)✓ Linux arm$(NC)"

# macOS builds
mac:
	@echo "$(BLUE)=== macOS Builds ===$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	@echo "$(GREEN)✓ macOS amd64 (Intel)$(NC)"
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	@echo "$(GREEN)✓ macOS arm64 (Apple Silicon)$(NC)"

# ARM builds
arm:
	@echo "$(BLUE)=== ARM Builds ===$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm .
	@echo "$(GREEN)✓ ARM$(NC)"
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .
	@echo "$(GREEN)✓ ARM64$(NC)"

# Run the application
run:
	@echo "$(BLUE)Running application...$(NC)"
	$(GOCMD) run main.go

# Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

# Build Docker image
docker:
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build -t $(APP_NAME):$(VERSION) .
	docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "$(GREEN)✓ Docker image built$(NC)"

# Create release archives
release: build-all
	@echo "$(BLUE)Creating release archives...$(NC)"
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && \
	for file in $(APP_NAME)-*; do \
		if [ -f "$$file" ]; then \
			platform=$$(echo $$file | sed 's/$(APP_NAME)-//'); \
			echo "Packaging $$platform..."; \
			tar -czf release/$(APP_NAME)-$$platform-$(VERSION).tar.gz $$file; \
		fi; \
	done
	@echo "$(GREEN)✓ Release archives created in $(BUILD_DIR)/release/$(NC)"

# Development mode with auto-reload (requires air)
dev:
	@echo "$(BLUE)Starting development server with auto-reload...$(NC)"
	@which air > /dev/null || (echo "$(RED)air not installed. Run: go install github.com/cosmtrek/air@latest$(NC)" && exit 1)
	air

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOCMD) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# Lint code
lint:
	@echo "$(BLUE)Linting code...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	golangci-lint run ./...
	@echo "$(GREEN)✓ Linting complete$(NC)"

# Show build info
info:
	@echo "$(BLUE)Build Information:$(NC)"
	@echo "  App Name:    $(APP_NAME)"
	@echo "  Version:     $(VERSION)"
	@echo "  Build Dir:   $(BUILD_DIR)"
	@echo "  Go Version:  $$(go version)"
	@echo ""

# Generate JWT secret
gen-secret:
	@echo "$(BLUE)Generating JWT Secret...$(NC)"
	@go run scripts/generate-jwt-secret.go

# Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"
