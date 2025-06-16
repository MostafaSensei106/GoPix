#!/bin/bash

# GoPix Installer for Linux, macOS
# Author: Mr. Mostafa Sensei
# Version: 1.0.0

set -euo pipefail

# Configuration
readonly APP_NAME="GoPix"
readonly BIN_NAME="GoPix"
readonly INSTALL_DIR="$HOME/.local/bin"
readonly VERSION="1.0.0"

# Colors
readonly GREEN='\033[0;32m'
readonly RED='\033[0;31m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Functions
show_help() {
    cat << EOF
${GREEN}Usage:${NC}
    $0 [options]

${GREEN}Options:${NC}
    -h, --help      Show this help message
    -r, --remove    Remove $BIN_NAME from $INSTALL_DIR
    -f, --force     Force reinstallation if already installed
    -v, --version   Show version information

${GREEN}Description:${NC}
    This script installs $APP_NAME, a Go-based image processing tool.
    
${GREEN}Requirements:${NC}
    - Go (golang) compiler
    - Git version control system
    - Unix-like operating system (Linux/macOS)

${GREEN}Examples:${NC}
    $0              # Install $APP_NAME
    $0 -r           # Remove $APP_NAME
    $0 -f           # Force reinstall
EOF
}

show_version() {
    echo -e "${BLUE}$APP_NAME Installer${NC}"
    echo -e "${BLUE}Version: $VERSION${NC}"
    echo -e "${BLUE}Author: Mr. Mostafa Sensei${NC}"
}

log_info() {
    echo -e "${GREEN}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}" >&2
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check system dependencies
check_dependency() {
    local dep="$1"
    if ! command_exists "$dep"; then
        log_error "Missing dependency: $dep. Please install it and try again."
        case "$dep" in
            "go")
                echo "  Install Go from: https://golang.org/dl/"
                ;;
            "git")
                echo "  Install Git from: https://git-scm.com/downloads"
                ;;
        esac
        exit 1
    fi
    log_success "$dep is installed."
}

# Detect OS and architecture
detect_system() {
    local os_name
    os_name=$(uname -s)
    
    log_info "Detected OS: $os_name"
    
    case "$os_name" in
        "Linux"|"Darwin")
            log_success "Your OS is supported"
            ;;
        *)
            log_error "Unsupported OS: $os_name"
            echo "This script only supports Linux and macOS (Darwin)"
            exit 1
            ;;
    esac
}

# Remove existing installation
remove_installation() {
    log_info "Uninstalling $APP_NAME from $INSTALL_DIR..."
    
    if [[ -f "$INSTALL_DIR/$BIN_NAME" ]]; then
        if rm -f "$INSTALL_DIR/$BIN_NAME"; then
            log_success "Removed $BIN_NAME successfully"
        else
            log_error "Failed to remove $BIN_NAME"
            exit 1
        fi
    else
        log_warning "$BIN_NAME is not installed in $INSTALL_DIR"
    fi
}

# Check if already installed
check_existing_installation() {
    if [[ -f "$INSTALL_DIR/$BIN_NAME" ]] && [[ "${FORCE_INSTALL:-}" != "true" ]]; then
        log_success "$BIN_NAME is already installed in $INSTALL_DIR"
        echo "Use -f or --force to reinstall"
        exit 0
    fi
}

# Build the application
build_application() {
    log_info "Building $APP_NAME..."
    
    # Check if we're in a Go project directory
    if [[ ! -f "go.mod" ]]; then
        log_error "go.mod file not found. Please run this script from the project root directory."
        exit 1
    fi
    
    # Build with proper flags
    local ldflags="-X 'github.com/mostafasensei106/gopix/cmd.Version=$VERSION' -s -w"
    
    if go build -x -ldflags "$ldflags" -o "$BIN_NAME" .; then
        log_success "$BIN_NAME built successfully!"
    else
        log_error "Failed to build $BIN_NAME"
        exit 1
    fi
}

# Install the binary
install_binary() {
    log_info "Installing to $INSTALL_DIR..."
    
    # Create install directory if it doesn't exist
    if ! mkdir -p "$INSTALL_DIR"; then
        log_error "Failed to create directory $INSTALL_DIR"
        exit 1
    fi
    
    # Move binary to install directory
    if ! mv "$BIN_NAME" "$INSTALL_DIR/"; then
        log_error "Failed to install $BIN_NAME to $INSTALL_DIR"
        exit 1
    fi
    
    # Make sure it's executable
    chmod +x "$INSTALL_DIR/$BIN_NAME"
    
    log_success "Binary installed successfully"
}

# Add to PATH if needed
update_path() {
    if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
        log_info "$INSTALL_DIR is already in PATH"
        return 0
    fi
    
    local shell_rc=""
    local path_command=""
    
    # Determine shell and appropriate RC file
    case "${SHELL##*/}" in
        "bash")
            shell_rc="$HOME/.bashrc"
            path_command="export PATH=\"\$PATH:$INSTALL_DIR\""
            ;;
        "zsh")
            shell_rc="$HOME/.zshrc"
            path_command="export PATH=\"\$PATH:$INSTALL_DIR\""
            ;;
        "fish")
            shell_rc="$HOME/.config/fish/config.fish"
            path_command="set -x PATH $INSTALL_DIR \$PATH"
            # Create fish config directory if it doesn't exist
            mkdir -p "$(dirname "$shell_rc")"
            ;;
        *)
            shell_rc="$HOME/.profile"
            path_command="export PATH=\"\$PATH:$INSTALL_DIR\""
            ;;
    esac
    
    # Check if PATH is already added (avoid duplicates)
    if [[ -f "$shell_rc" ]] && grep -q "$INSTALL_DIR" "$shell_rc" 2>/dev/null; then
        log_info "PATH already configured in $shell_rc"
        return 0
    fi
    
    # Add to PATH
    if echo "$path_command" >> "$shell_rc"; then
        log_success "Added $INSTALL_DIR to PATH in $shell_rc"
        log_warning "Please restart your terminal or run: source $shell_rc"
    else
        log_warning "Failed to update PATH automatically. Please add $INSTALL_DIR to your PATH manually."
    fi
}

# Confirm installation
confirm_installation() {
    log_info "👋 Hi there! I'm Mr. Mostafa Sensei, and this script will install $APP_NAME for you."
    
    while true; do
        read -rp "Continue with installation? (y/n): " answer
        case "$answer" in
            [Yy]|[Yy][Ee][Ss])
                break
                ;;
            [Nn]|[Nn][Oo])
                log_warning "Installation cancelled by user."
                exit 0
                ;;
            *)
                echo "Please answer yes (y) or no (n)."
                ;;
        esac
    done
}

# Main installation process
main_install() {
    log_info "🔍 Checking system requirements..."
    detect_system
    check_dependency "go"
    check_dependency "git"
    
    confirm_installation
    check_existing_installation
    
    build_application
    install_binary
    update_path
    
    log_success "🎉 Installation complete!"
    echo
    log_info "Try running: $BIN_NAME --help"
    log_info "Have a nice day!"
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                show_version
                exit 0
                ;;
            -r|--remove)
                remove_installation
                exit 0
                ;;
            -f|--force)
                FORCE_INSTALL="true"
                ;;
            -*)
                log_error "Unknown option: $1"
                echo "Use -h or --help for usage."
                exit 1
                ;;
            *)
                log_error "Unknown argument: $1"
                echo "Use -h or --help for usage."
                exit 1
                ;;
        esac
        shift
    done
}

# Main execution
main() {
    # Parse arguments first
    parse_arguments "$@"
    
    # If we reach here, proceed with installation
    main_install
}

# Run main function with all arguments
main "$@"