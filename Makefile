# Makefile for GoPix Installer (Cross-platform Enhanced)
# Author: Mr. Mostafa Sensei
# Version: 2.0.0

APP_NAME := GoPix
VERSION := 1.5.0
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
INSTALL_DIR := $(HOME)/.local/bin
CONFIG_DIR := $(HOME)/.config/$(APP_NAME)
LOG_DIR := $(HOME)/.local/share/$(APP_NAME)/logs

# Auto-detect OS and ARCH
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Enhanced OS Detection
ifeq ($(UNAME_S),Linux)
    GO_OS := linux
    INSTALL_DIR := $(HOME)/.local/bin
    CONFIG_DIR := $(HOME)/.config/$(APP_NAME)
else ifeq ($(UNAME_S),Darwin)
    GO_OS := darwin
    INSTALL_DIR := /usr/local/bin
    CONFIG_DIR := $(HOME)/Library/Application\ Support/$(APP_NAME)
else
    GO_OS := unknown
endif

# Enhanced Architecture Detection
ifeq ($(UNAME_M),x86_64)
    GO_ARCH := amd64
else ifeq ($(UNAME_M),aarch64)
    GO_ARCH := arm64
else ifeq ($(UNAME_M),arm64)
    GO_ARCH := arm64
else ifeq ($(UNAME_M),armv7l)
    GO_ARCH := arm
else ifeq ($(UNAME_M),i386)
    GO_ARCH := 386
else ifeq ($(UNAME_M),i686)
    GO_ARCH := 386
else
    GO_ARCH := unknown
endif

# Binary name based on OS
ifeq ($(GO_OS),windows)
    BIN_NAME := $(APP_NAME).exe
else
    BIN_NAME := $(APP_NAME)
endif

# Build flags with enhanced metadata
GO_BUILD_FLAGS := -ldflags "-X 'github.com/mostafasensei106/gopix/cmd.Version=$(VERSION)' \
                            -X 'github.com/mostafasensei106/gopix/cmd.BuildTime=$(BUILD_TIME)' \
                            -X 'github.com/mostafasensei106/gopix/cmd.GitCommit=$(GIT_COMMIT)' \
                            -s -w"

# Build modes
BUILD_MODE ?= release
ifeq ($(BUILD_MODE),debug)
    GO_BUILD_FLAGS := -race -ldflags "-X 'github.com/mostafasensei106/gopix/cmd.Version=$(VERSION)-debug' \
                                    -X 'github.com/mostafasensei106/gopix/cmd.BuildTime=$(BUILD_TIME)' \
                                    -X 'github.com/mostafasensei106/gopix/cmd.GitCommit=$(GIT_COMMIT)'"
endif

.PHONY: help install uninstall version build check-deps force-install update \
        clean test lint format dev-setup backup restore status doctor \
        build-all cross-compile

help: ## Show this help message
	@echo "$(APP_NAME) Installer v1.5.0"
	@echo "Author: Mr. Mostafa Sensei"
	@echo ""
	@echo "Usage: make [target] [BUILD_MODE=debug|release]"
	@echo ""
	@echo "Main Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "  %-20s %s\n", "Target", "Description"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Environment Variables:"
	@echo "  BUILD_MODE=debug     Build with debug symbols and race detection"
	@echo "  BUILD_MODE=release   Build optimized binary (default)"
	@echo "  INSTALL_DIR=path     Custom installation directory"

version: ## Show installer and app version
	@echo "$(APP_NAME) Installer"
	@echo "Version: $(VERSION)"
	@echo "Author: Mr. Mostafa Sensei"
	@echo "Platform: $(GO_OS)/$(GO_ARCH)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Mode: $(BUILD_MODE)"

status: ## Show current installation status
	@echo "🔍 Checking $(APP_NAME) status..."
	@echo "Target Platform: $(GO_OS)/$(GO_ARCH)"
	@echo "Install Directory: $(INSTALL_DIR)"
	@echo "Config Directory: $(CONFIG_DIR)"
	@if [ -f $(INSTALL_DIR)/$(BIN_NAME) ]; then \
		echo "✅ $(APP_NAME) is installed"; \
		echo "Version: $$($(INSTALL_DIR)/$(BIN_NAME) --version 2>/dev/null || echo 'Unknown')"; \
		echo "Size: $$(du -h $(INSTALL_DIR)/$(BIN_NAME) | cut -f1)"; \
		echo "Modified: $$(stat -c %y $(INSTALL_DIR)/$(BIN_NAME) 2>/dev/null || stat -f %Sm $(INSTALL_DIR)/$(BIN_NAME))"; \
	else \
		echo "❌ $(APP_NAME) is not installed"; \
	fi

doctor: ## Run system diagnostics
	@echo "🏥 Running system diagnostics..."
	@echo ""
	@echo "System Information:"
	@echo "  OS: $(UNAME_S)"
	@echo "  Architecture: $(UNAME_M)"
	@echo "  Go Target: $(GO_OS)/$(GO_ARCH)"
	@echo ""
	@echo "Dependencies:"
	@command -v go >/dev/null 2>&1 && echo "  ✅ Go: $$(go version)" || echo "  ❌ Go: Not installed"
	@command -v git >/dev/null 2>&1 && echo "  ✅ Git: $$(git --version)" || echo "  ❌ Git: Not installed"
	@command -v make >/dev/null 2>&1 && echo "  ✅ Make: $$(make --version | head -1)" || echo "  ❌ Make: Not installed"
	@echo ""
	@echo "Directories:"
	@echo "  Install: $(INSTALL_DIR) $(if $(wildcard $(INSTALL_DIR)),✅,❌)"
	@echo "  Config: $(CONFIG_DIR) $(if $(wildcard $(CONFIG_DIR)),✅,❌)"
	@echo ""
	@echo "PATH Check:"
	@case :$${PATH}: in *:$(INSTALL_DIR):*) echo "  ✅ Install directory is in PATH" ;; *) echo "  ⚠️  Install directory not in PATH" ;; esac

check-deps: ## Check system dependencies
	@echo "🔍 Checking dependencies..."
	@command -v go >/dev/null 2>&1 || { echo '❌ Go is not installed. Visit: https://golang.org/dl/'; exit 1; }
	@echo "✅ Go is installed: $$(go version)"
	@command -v git >/dev/null 2>&1 || { echo '❌ Git is not installed. Visit: https://git-scm.com/downloads'; exit 1; }
	@echo "✅ Git is installed: $$(git --version)"
	@if [ "$(GO_OS)" = "unknown" ] || [ "$(GO_ARCH)" = "unknown" ]; then \
		echo "❌ Unsupported platform: $(UNAME_S)/$(UNAME_M)"; \
		exit 1; \
	fi
	@echo "✅ Platform supported: $(GO_OS)/$(GO_ARCH)"

test: ## Run tests
	@echo "🧪 Running tests..."
	@test -f go.mod || { echo "❌ go.mod not found. Run from project root."; exit 1; }
	go test -v ./...
	@echo "✅ Tests completed."

lint: ## Run linting
	@echo "🔍 Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "⚠️  golangci-lint not installed. Skipping..."; exit 0; }
	golangci-lint run
	@echo "✅ Linting completed."

format: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted."

build: check-deps ## Build the Go binary
	@echo "🔧 Building $(APP_NAME) for $(GO_OS)/$(GO_ARCH) ($(BUILD_MODE) mode)..."
	@test -f go.mod || { echo "❌ go.mod not found. Run from project root."; exit 1; }
	@mkdir -p bin
	GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build $(GO_BUILD_FLAGS) -o bin/$(BIN_NAME) .
	@echo "✅ Built bin/$(BIN_NAME) successfully."
	@echo "📊 Binary size: $$(du -h bin/$(BIN_NAME) | cut -f1)"

build-all: ## Build for all supported platforms
	@echo "🌍 Building for all platforms..."
	@mkdir -p dist
	@for os in linux darwin windows freebsd; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then ext=".exe"; else ext=""; fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build $(GO_BUILD_FLAGS) -o dist/$(APP_NAME)-$$os-$$arch$$ext .; \
		done; \
	done
	@echo "✅ All builds completed in dist/ directory."

cross-compile: build-all ## Alias for build-all

dev-setup: check-deps ## Setup development environment
	@echo "🛠️  Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✅ Development environment ready."

install: check-deps build ## Install GoPix
	@echo "📦 Installing $(APP_NAME)..."
	@mkdir -p $(INSTALL_DIR)
	@mkdir -p $(CONFIG_DIR)
	@if [ -f $(INSTALL_DIR)/$(BIN_NAME) ]; then \
		echo "ℹ️  Backing up existing installation..."; \
		cp $(INSTALL_DIR)/$(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME).backup; \
	fi
	@cp bin/$(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)
	@chmod +x $(INSTALL_DIR)/$(BIN_NAME)
	@echo "✅ Installed to $(INSTALL_DIR)/$(BIN_NAME)"
	@echo "📁 Config directory: $(CONFIG_DIR)"
	@case :$${PATH}: in \
		*:$(INSTALL_DIR):*) echo "ℹ️  $(INSTALL_DIR) is in PATH" ;; \
		*) echo "⚠️  Add $(INSTALL_DIR) to your PATH:"; \
		   echo "   export PATH=\"$(INSTALL_DIR):\$$PATH\"" ;; \
	esac

force-install: ## Force reinstall GoPix
	@echo "🔄 Force reinstalling $(APP_NAME)..."
	$(MAKE) uninstall || true
	$(MAKE) install

backup: ## Backup current installation
	@echo "💾 Creating backup..."
	@if [ -f $(INSTALL_DIR)/$(BIN_NAME) ]; then \
		backup_name="$(APP_NAME)-backup-$$(date +%Y%m%d-%H%M%S)"; \
		cp $(INSTALL_DIR)/$(BIN_NAME) $(INSTALL_DIR)/$$backup_name; \
		echo "✅ Backup created: $(INSTALL_DIR)/$$backup_name"; \
	else \
		echo "⚠️  No installation found to backup"; \
	fi

restore: ## Restore from backup
	@echo "🔄 Restoring from backup..."
	@backup_file=$$(ls -t $(INSTALL_DIR)/$(APP_NAME)-backup-* 2>/dev/null | head -1); \
	if [ -n "$$backup_file" ]; then \
		cp "$$backup_file" $(INSTALL_DIR)/$(BIN_NAME); \
		chmod +x $(INSTALL_DIR)/$(BIN_NAME); \
		echo "✅ Restored from: $$backup_file"; \
	else \
		echo "❌ No backup found"; \
		exit 1; \
	fi

update: ## Update GoPix from git
	@echo "🔄 Updating $(APP_NAME)..."
	@if [ -d .git ]; then \
		echo "📡 Pulling latest changes..."; \
		git pull --rebase || { echo "⚠️  Git pull failed"; exit 1; }; \
		$(MAKE) backup; \
		$(MAKE) force-install; \
		echo "✅ Update completed"; \
	else \
		echo "❌ Not in a git repository. Cannot update."; \
		echo "💡 Try downloading the latest version manually."; \
		exit 1; \
	fi

uninstall: ## Remove GoPix
	@echo "🗑️  Uninstalling $(APP_NAME)..."
	@if [ -f $(INSTALL_DIR)/$(BIN_NAME) ]; then \
		rm -f $(INSTALL_DIR)/$(BIN_NAME); \
		echo "✅ Removed $(BIN_NAME) from $(INSTALL_DIR)"; \
	else \
		echo "⚠️  $(BIN_NAME) not found in $(INSTALL_DIR)"; \
	fi
	@echo "ℹ️  Config directory preserved: $(CONFIG_DIR)"
	@echo "💡 To remove config: rm -rf $(CONFIG_DIR)"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ dist/ $(BIN_NAME)
	@echo "✅ Cleaned build artifacts."

clean-all: clean ## Clean everything including config
	@echo "🧹 Cleaning everything..."
	@read -p "Remove config directory $(CONFIG_DIR)? [y/N]: " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		rm -rf $(CONFIG_DIR); \
		echo "✅ Removed config directory"; \
	fi

# Advanced targets
install-global: check-deps build ## Install globally (requires sudo)
	@echo "🌍 Installing $(APP_NAME) globally..."
	@if [ "$(GO_OS)" = "darwin" ]; then \
		sudo cp bin/$(BIN_NAME) /usr/local/bin/$(BIN_NAME); \
	else \
		sudo cp bin/$(BIN_NAME) /usr/bin/$(BIN_NAME); \
	fi
	@echo "✅ Installed globally"

install-user: install ## Install for current user (default)

# Development targets
watch: ## Watch and rebuild on changes (requires entr)
	@command -v entr >/dev/null 2>&1 || { echo "❌ entr not installed. Install with: apt install entr"; exit 1; }
	@echo "👀 Watching for changes..."
	find . -name "*.go" | entr -r make build

# Release targets
release: ## Create a release build
	@$(MAKE) clean
	@$(MAKE) test
	@$(MAKE) build-all
	@echo "✅ Release builds ready in dist/"

# Quick development workflow
dev: dev-setup format lint test build ## Full development workflow

.DEFAULT_GOAL := help