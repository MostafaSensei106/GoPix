#!/bin/bash

set -e

APP_NAME="imgconvert"
BIN_NAME="imgconvert"
INSTALL_DIR="$HOME/.local/bin"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

check_dependency() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}❌ Missing dependency: $1. Please install it and try again.${NC}"
        exit 1
    fi
}

echo -e "${GREEN}🔍 Checking system requirements...${NC}"
check_dependency "go"

platform=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

case $platform in
    linux*|darwin*)
        ;;
    *)
        echo -e "${RED}❌ Unsupported platform: $platform${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}👋 Hi There Iam Mr.Mostafa Sensei! And This Script Will Install ${APP_NAME}...${NC}"

read -p "Continue with installation? (y/n): " answer
if [[ "$answer" != "y" && "$answer" != "Y" ]]; then
    echo -e "${RED}❌ Installation cancelled.${NC}"
    exit 0
fi

echo -e "${GREEN}🔧 Building $APP_NAME...${NC}"
go build -ldflags "-X 'github.com/mostafasensei106/gopix/cmd.Version=1.0.0'" -o "$BIN_NAME"

echo -e "${GREEN}📦 Installing to $INSTALL_DIR...${NC}"
mkdir -p "$INSTALL_DIR"
mv -f "$BIN_NAME" "$INSTALL_DIR/"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    SHELL_RC="$HOME/.bashrc"
    [[ $SHELL == */zsh ]] && SHELL_RC="$HOME/.zshrc"
    [[ $SHELL == */fish ]] && SHELL_RC="$HOME/.config/fish/config.fish"

    if [[ $SHELL == */fish ]]; then
        echo "set -x PATH $INSTALL_DIR \$PATH" >> "$SHELL_RC"
    else
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
    fi

    echo -e "${GREEN}📌 Added $INSTALL_DIR to PATH in $SHELL_RC (restart terminal to apply)${NC}"
fi

echo -e "${GREEN}🎉 Installation complete! Try running:${NC} $BIN_NAME --help"
