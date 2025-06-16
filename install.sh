#!/bin/bash

set -e

APP_NAME="GoPix"
BIN_NAME="GoPix"
INSTALL_DIR="$HOME/.local/bin"
os_name=$(uname -s)

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  echo -e "${GREEN}Usage:${NC}"
  echo -e "${GREEN}  $0 [options]${NC}"
  echo -e "${GREEN}Options:${NC}"
  echo -e "${GREEN}  -h, --help${NC}   Show this help message"
  echo -e "${GREEN}  -r, --remove${NC} Remove $BIN_NAME from $INSTALL_DIR"
  exit 0
else
  echo "Unknown option or no option provided. Use -h or --help for usage."
  exit 1
fi

echo -e "${GREEN}🖥️  Detected OS Is $os_name...${NC}"
if [[ "$os_name" != "Linux" && "$os_name" != "Darwin" ]]; then
  echo -e "${RED}❌ Unsupported OS: $os_name${NC}"
  exit 1
fi
echo -e "${GREEN}✅ Your OS is supported${NC}"

if [[ "$1" == "-r" || "$1" == "--remove" ]]; then
  echo "🧹 Uninstalling $APP_NAME from $INSTALL_DIR ..."
  if [[ -f "$INSTALL_DIR/$BIN_NAME" ]]; then
    rm -f "$INSTALL_DIR/$BIN_NAME"
    echo "✅ Removed $BIN_NAME"
  else
    echo "⚠️  $BIN_NAME is not installed in $INSTALL_DIR"
  fi
  exit 0
fi

if [[ -f "$INSTALL_DIR/$BIN_NAME" ]]; then
  echo -e "${GREEN}✅ $BIN_NAME is already installed in $INSTALL_DIR${NC}"
  exit 0
fi



check_dependency() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}❌ Missing dependency: $1. Please install it and try again.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ $1 is installed.${NC}"
}



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

echo -e "${GREEN}👋 Hi there! I'm Mr. Mostafa Sensei, and this script will install ${APP_NAME} for you.${NC}"

read -p "Continue with installation? (y/n): " answer
if [[ "$answer" != "y" && "$answer" != "Y" ]]; then
    echo -e "${RED}❌ Installation cancelled.${NC}"
    exit 0
fi

echo -e "${GREEN}🔍 Checking system requirements...${NC}"
check_dependency "go"

echo -e "${GREEN}🔧 Building $APP_NAME...${NC}"
go build -ldflags "-X 'github.com/mostafasensei106/gopix/cmd.Version=1.0.0'" -o "$BIN_NAME"
echo -e "${GREEN}✅ $BIN_NAME built successfully!${NC}"

echo -e "${GREEN}📦 Installing to $INSTALL_DIR...${NC}"
mkdir -p "$INSTALL_DIR"
mv -f "$BIN_NAME" "$INSTALL_DIR/"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  SHELL_RC=""
  case "$SHELL" in
    */bash) SHELL_RC="$HOME/.bashrc" ;;
    */zsh)  SHELL_RC="$HOME/.zshrc" ;;
    */fish) SHELL_RC="$HOME/.config/fish/config.fish" ;;
    *) SHELL_RC="$HOME/.profile" ;;  # fallback
  esac

  if [[ "$SHELL" == */fish ]]; then
    # Avoid duplicate
    if ! grep -q "set -x PATH $INSTALL_DIR" "$SHELL_RC" 2>/dev/null; then
      echo "set -x PATH $INSTALL_DIR \$PATH" >> "$SHELL_RC"
      echo -e "${GREEN}📌 Added $INSTALL_DIR to PATH in $SHELL_RC (fish)${NC}"
    fi
  else
    # Avoid duplicate
    if ! grep -q "export PATH=.*$INSTALL_DIR" "$SHELL_RC" 2>/dev/null; then
      echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
      echo -e "${GREEN}📌 Added $INSTALL_DIR to PATH in $SHELL_RC${NC}"
    fi
  fi

  echo -e "${YELLOW}🔁 Please restart your terminal to apply changes.${NC}"
fi

echo -e "${GREEN}🎉 Installation complete! Try running:${NC} $BIN_NAME --help"
echo -e "${GREEN}Have a nice day!${NC}"
