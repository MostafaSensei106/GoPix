# Makefile for building and installing GoPix
# Author: Mostafa Sensei106
# License: Global Public License v3 (GPL v3)

# declare variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
APP_NAME := gopix
OUTPUT_DIR := bin/$(GOOS)/$(GOARCH)
OUTPUT := $(OUTPUT_DIR)/$(APP_NAME)
GoPix_VERSION := 1.5.2

# declare installation directories
INSTALL_DIR_LINUX := /usr/local/bin
INSTALL_DIR_WIN := /c/Program\ Files/$(APP_NAME)/bin

.PHONY: all build install clean release help check

all: check build

deps:
	@echo "📦 Checking dependencies..."
	@if [ -f go.sum ]; then \
		echo "📦 Verifying dependencies..."; \
		go mod verify; \
		echo "✅ Dependencies installed and up-to-date"; \
	else \
		echo "📦 Downloading dependencies..."; \
		go mod download; \
		echo "📦 Verifying dependencies..."; \
		go mod verify; \
		echo "✅ Dependencies installed"; \
	fi


check: deps


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


release: check
	@{ \
		echo "🔍 Detecting host platform..."; \
		HOST_OS=$$(go env GOOS); \
		HOST_ARCH=$$(go env GOARCH); \
		echo "🖥️  Host: $$HOST_OS/$$HOST_ARCH"; \
		echo "🌐  Building for all major platforms and architectures..."; \
		platforms="linux/amd64 windows/amd64"; \
		for platform in $$platforms; do \
			GOOS=$${platform%/*}; \
			GOARCH=$${platform#*/}; \
			OUT_DIR=bin/$$GOOS/$$GOARCH; \
			OUT_FILE=$$OUT_DIR/$(APP_NAME); \
			if [ "$$GOOS" = "windows" ]; then \
				OUT_FILE=$$OUT_FILE.exe; \
			fi; \
			ARCHIVE_NAME=$(APP_NAME)-$$GOOS-$$GOARCH-${GoPix_VERSION}; \
			mkdir -p $$OUT_DIR; \
			echo "🛠️  Building for $$GOOS/$$GOARCH..."; \
			if [ "$$GOOS" = "windows" ]; then \
				if [ "$$HOST_OS" = "windows" ]; then \
					GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE || echo "❌ Failed for $$GOOS/$$GOARCH"; \
				elif command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1; then \
					CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE || echo "❌ Failed for $$GOOS/$$GOARCH"; \
				else \
					echo "⚠️ Skipped: $$GOOS/$$GOARCH (missing cross-compiler)"; \
					continue; \
				fi; \
			elif [ "$$GOOS" = "darwin" ] && [ "$$HOST_OS" != "darwin" ]; then \
				echo "⚠️ Skipped: $$GOOS/$$GOARCH (macOS cross-compilation unsupported outside macOS)"; \
				continue; \
			else \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUT_FILE || echo "❌ Failed for $$GOOS/$$GOARCH"; \
			fi; \
			echo "✅ Done: $$OUT_FILE"; \
			mkdir -p release; \
			if [ "$$GOOS" = "windows" ]; then \
				cd bin && zip -r "../release/$$ARCHIVE_NAME.zip" "$$GOOS/$$GOARCH" >/dev/null && cd .. && \
				echo "✅ Compressed (zip): release/$$ARCHIVE_NAME.zip"; \
			else \
				cd bin && tar -czf "../release/$$ARCHIVE_NAME.tar.gz" "$$GOOS/$$GOARCH" >/dev/null && cd .. && \
				echo "✅ Compressed (tar.gz): release/$$ARCHIVE_NAME.tar.gz"; \
			fi; \
		done; \
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
	@echo "make check    👉 Run all checks (deps, fmt, vet)"
	@echo "make build    👉 Build the app for current OS/arch"
	@echo "make install  👉 Build and install GoPix to system"
	@echo "make release  👉 Build for all OS/platforms for release"
	@echo "make clean    👉 Delete all build artifacts and caches"
	@echo "make help     👉 Show this help message"
	@echo ""
