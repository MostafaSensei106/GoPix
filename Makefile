# Makefile for GoPix Installer
# Author: Mr. Mostafa Sensei
# Version: 1.5.0 

APP_NAME := GoPix
VERSION := 1.5.0
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Auto-detect OS and ARCH
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# OS Detection and default paths
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

# Architecture Detection
ifeq ($(UNAME_M),x86_64)
    GO_ARCH := amd64
else ifeq ($(UNAME_M),aarch64)
    GO_ARCH := arm64
else ifeq ($(UNAME_M),arm64) # macOS M1/M2
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


BIN_NAME := $(APP_NAME)

# Build flags with essential metadata
GO_BUILD_FLAGS := -ldflags "-X 'github.com/mostafasensei106/gopix/cmd.Version=$(VERSION)' \
                            -X 'github.com/mostafasensei106/gopix/cmd.BuildTime=$(BUILD_TIME)' \
                            -s -w"

# Define phony targets for make
.PHONY: help build install uninstall clean check-deps status

help: ## Show this help message
    @echo "$(APP_NAME) Installer v$(VERSION)"
    @echo "Author: Mr. Mostafa Sensei"
    @echo ""
    @echo "Usage: make [target]"
    @echo ""
    @echo "Available Targets:"
    @awk 'BEGIN {FS = ":.*##"; printf "  %-15s %s\n", "Target", "Description"} \
        /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

check-deps: ## Check system dependencies (Go and Git)
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

build: check-deps ## Build the Go binary for your OS
    @echo "🔧 Building $(APP_NAME) for $(GO_OS)/$(GO_ARCH)..."
    @test -f go.mod || { echo "❌ go.mod not found. Run from project root."; exit 1; }
    @mkdir -p bin
    GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build $(GO_BUILD_FLAGS) -o bin/$(BIN_NAME) .
    @echo "✅ Built bin/$(BIN_NAME) successfully."
    @echo "📊 Binary size: $$(du -h bin/$(BIN_NAME) | cut -f1)"

install: check-deps build ## Install GoPix to your system
    @echo "📦 Installing $(APP_NAME)..."
    # Special check for macOS
    ifeq ($(GO_OS),darwin)
        @echo "😂 هذه الأدوات ليست للأغبياء! لا يمكن تثبيت GoPix على نظام macOS باستخدام هذا Makefile."
        @echo "💡 لو عايز تثبتها على ماك، ممكن تاخد الملف التنفيذي من فولدر 'bin/' وتنقله يدويًا."
        @exit 1
    endif
    # Proceed with installation for other OS (Linux)
    @mkdir -p $(INSTALL_DIR)
    @mkdir -p $(CONFIG_DIR) # Create config directory for all OS
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

uninstall: ## Remove GoPix from your system
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
    @rm -rf bin/
    @echo "✅ Cleaned build artifacts."

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

.DEFAULT_GOAL := help
