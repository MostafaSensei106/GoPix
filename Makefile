# Makefile for building and installing GoPix
# Author: Mostafa Sensei106
# License: Global Public License v3 (GPL v3)

# declare variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
APP_NAME := gopix
SRC := ./main.go
OUTPUT_DIR := bin/$(GOOS)/$(GOARCH)
OUTPUT := $(OUTPUT_DIR)/$(APP_NAME)

# declare installation directories
INSTALL_DIR_LINUX := $(HOME)/.local/bin
INSTALL_DIR_WIN := /c/Program\ Files/$(APP_NAME)/bin

.PHONY: all build install clean buildAll help deps check fmt vet

all: check build

deps:
	@echo "📦 Checking dependencies..."
	@if [ -f go.mod ]; then \
		echo "✅ Dependencies already downloaded"; \
	else \
		echo "📦 Downloading dependencies..."; \
		go mod download; \
		echo "✅ Dependencies downloaded"; \
	fi

fmt:
	@echo "🎨 Checking code formatting..."
	@if [ -n "$(gofmt -l .)" ]; then \
		echo "❌ Code is not formatted. Run 'go fmt ./...' to fix:"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "✅ Code is properly formatted"


vet:
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ go vet passed"


check: deps fmt vet
	@echo "🔍 Running all checks..."
	@echo "✅ All checks passed"

build: check
	@echo "📦 Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT) 
	@echo "✅ Build complete: $(OUTPUT)"

install: check build
ifeq ($(GOOS),windows)
	@echo "📥 Installing for Windows..."
	@mkdir -p $(INSTALL_DIR_WIN)
	@cp $(OUTPUT).exe $(INSTALL_DIR_WIN)/$(APP_NAME).exe
	@echo "✅ Installed to $(INSTALL_DIR_WIN)/$(APP_NAME).exe"
else
	@echo "📥 Installing for $$HOST_OS/$$HOST_ARCH system..."
	@mkdir -p $(INSTALL_DIR_LINUX)
	@sudo cp $(OUTPUT) $(INSTALL_DIR_LINUX)/$(APP_NAME)
	@echo "✅ Installed to $(INSTALL_DIR_LINUX)/$(APP_NAME)"
endif

buildAll: check
	@{ \
		echo "🔍 Detecting host platform..."; \
		HOST_OS=$$(go env GOOS); \
		HOST_ARCH=$$(go env GOARCH); \
		echo "🖥️  Host: $$HOST_OS/$$HOST_ARCH"; \
		echo "🌐  Building for all major platforms and architectures..."; \
		platforms="linux/amd64 windows/amd64 darwin/arm64"; \
		for platform in $$platforms; do \
			GOOS=$${platform%/*}; \
			GOARCH=$${platform#*/}; \
			OUT_DIR=bin/$$GOOS/$$GOARCH; \
			OUT_FILE=$$OUT_DIR/$(APP_NAME); \
			if [ "$$GOOS" = "windows" ]; then \
				OUT_FILE=$$OUT_FILE.exe; \
			fi; \
			mkdir -p $$OUT_DIR; \
			echo "	🛠️ Building for $$GOOS/$$GOARCH..."; \
			if [ "$$GOOS" = "windows" ]; then \
				if [ "$$HOST_OS" = "windows" ]; then \
					GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE $(SRC) && \
					echo "✅ Done: $$OUT_FILE" || echo "❌ Failed for $$GOOS/$$GOARCH"; \
				elif command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then \
					CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE $(SRC) && \
					echo "✅ Done: $$OUT_FILE" || echo "❌ Failed for $$GOOS/$$GOARCH"; \
				else \
					echo "⚠️ Skipped: $$GOOS/$$GOARCH (missing cross-compiler)"; \
				fi; \
			elif [ "$$GOOS" = "darwin" ] && [ "$$HOST_OS" != "darwin" ]; then \
				echo "⚠️ Skipped: $$GOOS/$$GOARCH (macOS cross-compilation unsupported outside macOS)"; \
			else \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE $(SRC) && \
				echo "✅ Done: $$OUT_FILE" || echo "✅ Failed for $$GOOS/$$GOARCH"; \
			fi; \
		done; \
		echo "✅ All builds attempted."; \
	}

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin
	@go clean -cache -modcache -testcache
	@echo "✅ Clean complete."

help:
	@echo ""
	@echo "📖 GoPix Makefile Commands"
	@echo "============================"
	@echo "make deps     👉 Download and update dependencies"
	@echo "make fmt      👉 Check code formatting"
	@echo "make vet      👉 Run go vet"
	@echo "make check    👉 Run all checks (deps, fmt, vet)"
	@echo "make build    👉 Build the app for current OS/arch"
	@echo "make install  👉 Build and install GoPix to system"
	@echo "make buildAll 👉 Build for all OS/platforms for release"
	@echo "make clean    👉 Delete all build artifacts and caches"
	@echo "make help     👉 Show this help message"
	@echo ""