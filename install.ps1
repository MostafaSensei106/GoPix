# PowerShell Installer for GoPix (Windows)

$AppName = "GoPix"
$BinaryName = "GoPix.exe"
$InstallDir = "$env:USERPROFILE\.bin"
$osName = [System.Runtime.InteropServices.RuntimeInformation]::OSDescription
$GoRequired = $true

Write-Host "🖥️ Detected OS: $osName"

if ($osName -notlike "*Windows*") {
    Write-Host "❌ Unsupported OS: $osName" -ForegroundColor Red
    exit 1
}


function Check-Dependency {
    param ([string]$cmd)
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Host "❌ Missing dependency: $cmd. Please install it first." -ForegroundColor Red
        exit 1
    }
}

Write-Host "🔍 Checking dependencies..." -ForegroundColor Cyan

if ($GoRequired) { Check-Dependency "go" }

Write-Host "✅ All dependencies are available." -ForegroundColor Green

$confirmation = Read-Host "Install $AppName to $InstallDir ? (Y/N)"
if ($confirmation -ne "Y" -and $confirmation -ne "y") {
    Write-Host "❌ Installation cancelled."
    exit 0
}

if (-Not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host "🔧 Building $AppName..." -ForegroundColor Cyan
go build -ldflags "-X 'imgconvert/cmd.Version=1.0.0'" -o $BinaryName

Write-Host "📦 Installing to $InstallDir..." -ForegroundColor Cyan
Move-Item -Force $BinaryName "$InstallDir\$BinaryName"

# Check if $InstallDir is in PATH
$pathList = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User).Split(";")
if ($InstallDir -notin $pathList) {
    Write-Host "🛠 Adding $InstallDir to user PATH..." -ForegroundColor Yellow
    $newPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + ";$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, [EnvironmentVariableTarget]::User)

    Write-Host "🔄 Please restart your terminal or log out/in to apply PATH changes." -ForegroundColor Yellow
}

Write-Host "`n🎉 $AppName installed successfully!"
Write-Host "👉 Run with: ${AppName} --help"
