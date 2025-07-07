# Makefile for building and installing GoPix

# Defaults (can be overridden via CLI)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP_NAME := gopix
SRC := ./main.go
OUTPUT_DIR := bin/$(GOOS)
OUTPUT := $(OUTPUT_DIR)/$(APP_NAME)

INSTALL_DIR_LINUX := ~/$USER/.local/bin
INSTALL_DIR_WIN := bin

.PHONY: all build install clean

all: build

build:
	@echo "📦 Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT) $(SRC)
	@echo "✅ Build complete: $(OUTPUT)"

install: build
ifeq ($(GOOS),windows)
	@echo "📥 Installing for Windows..."
	@mkdir -p $(INSTALL_DIR_WIN)
	cp $(OUTPUT) $(INSTALL_DIR_WIN)/$(APP_NAME).exe
	@echo "✅ Installed to $(INSTALL_DIR_WIN)/$(APP_NAME).exe"
else
	@echo "📥 Installing for Unix-like system..."
	sudo cp $(OUTPUT) $(INSTALL_DIR_LINUX)/$(APP_NAME)
	@echo "✅ Installed to $(INSTALL_DIR_LINUX)/$(APP_NAME)"
endif

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf build
	rm -rf bin
	@echo "✅ Clean complete."
