# GoPix Installer for Windows PowerShell
# Author: Mr. Mostafa Sensei
# Version: 1.0.0

param(
    [switch]$Help,
    [switch]$h,
    [switch]$Remove,
    [switch]$r,
    [switch]$Force,
    [switch]$f,
    [switch]$Version,
    [switch]$v
)

# Configuration
$Script:APP_NAME = "GoPix"
$Script:BIN_NAME = "GoPix.exe"
$Script:VERSION = "1.0.0"
$Script:INSTALL_DIR = "$env:USERPROFILE\.local\bin"

# Error handling
$ErrorActionPreference = "Stop"

# Color functions for better output
function Write-ColoredText {
    param(
        [string]$Text,
        [string]$Color = "White",
        [string]$Emoji = ""
    )
    
    $colorMap = @{
        "Green" = "Green"
        "Red" = "Red"
        "Yellow" = "Yellow"
        "Blue" = "Cyan"
        "White" = "White"
    }
    
    if ($Emoji) {
        Write-Host "$Emoji $Text" -ForegroundColor $colorMap[$Color]
    } else {
        Write-Host $Text -ForegroundColor $colorMap[$Color]
    }
}

function Write-Info {
    param([string]$Message)
    Write-ColoredText -Text $Message -Color "Green" -Emoji "ℹ️"
}

function Write-Success {
    param([string]$Message)
    Write-ColoredText -Text $Message -Color "Green" -Emoji "✅"
}

function Write-Error {
    param([string]$Message)
    Write-ColoredText -Text $Message -Color "Red" -Emoji "❌"
    Write-Error $Message
}

function Write-Warning {
    param([string]$Message)
    Write-ColoredText -Text $Message -Color "Yellow" -Emoji "⚠️"
}

# Help function
function Show-Help {
    Write-ColoredText -Text "Usage:" -Color "Green"
    Write-Host "    .\install.ps1 [options]"
    Write-Host ""
    Write-ColoredText -Text "Options:" -Color "Green"
    Write-Host "    -Help, -h       Show this help message"
    Write-Host "    -Remove, -r     Remove $Script:BIN_NAME from installation directory"
    Write-Host "    -Force, -f      Force reinstallation if already installed"
    Write-Host "    -Version, -v    Show version information"
    Write-Host ""
    Write-ColoredText -Text "Description:" -Color "Green"
    Write-Host "    This script installs $Script:APP_NAME, a Go-based image processing tool."
    Write-Host ""
    Write-ColoredText -Text "Requirements:" -Color "Green"
    Write-Host "    - Go (golang) compiler"
    Write-Host "    - Git version control system"
    Write-Host "    - Windows PowerShell 5.1 or PowerShell Core 6+"
    Write-Host ""
    Write-ColoredText -Text "Examples:" -Color "Green"
    Write-Host "    .\install.ps1           # Install $Script:APP_NAME"
    Write-Host "    .\install.ps1 -r        # Remove $Script:APP_NAME"
    Write-Host "    .\install.ps1 -f        # Force reinstall"
}

# Version function
function Show-Version {
    Write-ColoredText -Text "$Script:APP_NAME Installer" -Color "Blue"
    Write-ColoredText -Text "Version: $Script:VERSION" -Color "Blue"
    Write-ColoredText -Text "Author: Mr. Mostafa Sensei" -Color "Blue"
    Write-ColoredText -Text "Platform: Windows PowerShell" -Color "Blue"
}

# Check if command exists
function Test-CommandExists {
    param([string]$Command)
    
    try {
        Get-Command $Command -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

# Check system dependencies
function Test-Dependencies {
    Write-Info "🔍 Checking system requirements..."
    
    # Check PowerShell version
    $psVersion = $PSVersionTable.PSVersion
    Write-Info "PowerShell Version: $psVersion"
    
    if ($psVersion.Major -lt 5) {
        Write-Error "PowerShell 5.0 or higher is required. Current version: $psVersion"
        exit 1
    }
    
    # Check Go
    if (-not (Test-CommandExists "go")) {
        Write-Error "Missing dependency: Go. Please install it and try again."
        Write-Host "  Install Go from: https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
    Write-Success "Go is installed."
    
    # Check Git
    if (-not (Test-CommandExists "git")) {
        Write-Error "Missing dependency: Git. Please install it and try again."
        Write-Host "  Install Git from: https://git-scm.com/downloads" -ForegroundColor Yellow
        exit 1
    }
    Write-Success "Git is installed."
    
    Write-Success "All dependencies are satisfied."
}

# Remove existing installation
function Remove-Installation {
    Write-Info "Uninstalling $Script:APP_NAME from $Script:INSTALL_DIR..."
    
    $binaryPath = Join-Path $Script:INSTALL_DIR $Script:BIN_NAME
    
    if (Test-Path $binaryPath) {
        try {
            Remove-Item $binaryPath -Force
            Write-Success "Removed $Script:BIN_NAME successfully"
        }
        catch {
            Write-Error "Failed to remove $Script:BIN_NAME: $($_.Exception.Message)"
            exit 1
        }
    }
    else {
        Write-Warning "$Script:BIN_NAME is not installed in $Script:INSTALL_DIR"
    }
}

# Check if already installed
function Test-ExistingInstallation {
    $binaryPath = Join-Path $Script:INSTALL_DIR $Script:BIN_NAME
    
    if ((Test-Path $binaryPath) -and (-not $Script:FORCE_INSTALL)) {
        Write-Success "$Script:BIN_NAME is already installed in $Script:INSTALL_DIR"
        Write-Host "Use -Force or -f to reinstall"
        exit 0
    }
}

# Build the application
function Build-Application {
    Write-Info "Building $Script:APP_NAME..."
    
    # Check if we're in a Go project directory
    if (-not (Test-Path "go.mod")) {
        Write-Error "go.mod file not found. Please run this script from the project root directory."
        exit 1
    }
    
    # Build with proper flags
    $ldflags = "-X 'github.com/mostafasensei106/gopix/cmd.Version=$Script:VERSION' -s -w"
    
    try {
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        
        Write-Info "Building for Windows (amd64)..."
        & go build -x -ldflags $ldflags -o $Script:BIN_NAME .
        
        if ($LASTEXITCODE -ne 0) {
            throw "Go build failed with exit code $LASTEXITCODE"
        }
        
        Write-Success "$Script:BIN_NAME built successfully!"
    }
    catch {
        Write-Error "Failed to build $Script:BIN_NAME: $($_.Exception.Message)"
        exit 1
    }
}

# Install the binary
function Install-Binary {
    Write-Info "Installing to $Script:INSTALL_DIR..."
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $Script:INSTALL_DIR)) {
        try {
            New-Item -ItemType Directory -Path $Script:INSTALL_DIR -Force | Out-Null
            Write-Success "Created directory $Script:INSTALL_DIR"
        }
        catch {
            Write-Error "Failed to create directory $Script:INSTALL_DIR: $($_.Exception.Message)"
            exit 1
        }
    }
    
    # Move binary to install directory
    $sourcePath = Join-Path (Get-Location) $Script:BIN_NAME
    $destinationPath = Join-Path $Script:INSTALL_DIR $Script:BIN_NAME
    
    try {
        Move-Item $sourcePath $destinationPath -Force
        Write-Success "Binary installed successfully"
    }
    catch {
        Write-Error "Failed to install $Script:BIN_NAME to $Script:INSTALL_DIR: $($_.Exception.Message)"
        exit 1
    }
}

# Update PATH environment variable
function Update-PathEnvironment {
    # Check if install directory is already in PATH
    $currentPath = $env:PATH
    if ($currentPath -like "*$Script:INSTALL_DIR*") {
        Write-Info "$Script:INSTALL_DIR is already in PATH"
        return
    }
    
    try {
        # Get current user PATH
        $userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        
        if (-not $userPath) {
            $userPath = $Script:INSTALL_DIR
        }
        elseif ($userPath -notlike "*$Script:INSTALL_DIR*") {
            $userPath = "$userPath;$Script:INSTALL_DIR"
        }
        else {
            Write-Info "PATH already configured"
            return
        }
        
        # Set the new PATH
        [Environment]::SetEnvironmentVariable("PATH", $userPath, "User")
        
        # Update current session PATH
        $env:PATH = "$env:PATH;$Script:INSTALL_DIR"
        
        Write-Success "Added $Script:INSTALL_DIR to user PATH"
        Write-Warning "Please restart your terminal to apply PATH changes, or run: refreshenv"
    }
    catch {
        Write-Warning "Failed to update PATH automatically: $($_.Exception.Message)"
        Write-Warning "Please add $Script:INSTALL_DIR to your PATH manually."
        Write-Host "You can do this by:"
        Write-Host "1. Opening System Properties -> Environment Variables"
        Write-Host "2. Adding $Script:INSTALL_DIR to your user PATH variable"
    }
}

# Confirm installation
function Confirm-Installation {
    Write-Info "👋 Hi there! I'm Mr. Mostafa Sensei, and this script will install $Script:APP_NAME for you."
    
    do {
        $answer = Read-Host "Continue with installation? (y/n)"
        switch ($answer.ToLower()) {
            { $_ -in @("y", "yes") } {
                return
            }
            { $_ -in @("n", "no") } {
                Write-Warning "Installation cancelled by user."
                exit 0
            }
            default {
                Write-Host "Please answer yes (y) or no (n)."
            }
        }
    } while ($true)
}

# Main installation process
function Start-Installation {
    Test-Dependencies
    Confirm-Installation
    Test-ExistingInstallation
    Build-Application
    Install-Binary
    Update-PathEnvironment
    
    Write-Success "🎉 Installation complete!"
    Write-Host ""
    Write-Info "Try running: $($Script:BIN_NAME.Replace('.exe', '')) --help"
    Write-Info "Have a nice day!"
}

# Main execution logic
function Main {
    # Handle help
    if ($Help -or $h) {
        Show-Help
        exit 0
    }
    
    # Handle version
    if ($Version -or $v) {
        Show-Version
        exit 0
    }
    
    # Handle removal
    if ($Remove -or $r) {
        Remove-Installation
        exit 0
    }
    
    # Set force install flag
    if ($Force -or $f) {
        $Script:FORCE_INSTALL = $true
    }
    
    # Proceed with installation
    Start-Installation
}

# Entry point
try {
    Main
}
catch {
    Write-Error "An unexpected error occurred: $($_.Exception.Message)"
    Write-Host "Stack trace:" -ForegroundColor Red
    Write-Host $_.ScriptStackTrace -ForegroundColor Red
    exit 1
}