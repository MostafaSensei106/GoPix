package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var SupportedOutput = map[string]bool{
	"png":  true,
	"jpg":  true,
	"jpeg": true,
	"webp": true,
	"avif": true,
	"tiff": true,
	"bmp":  true,

	"heic": false,
	"heif": false,
	"ico":  false,
}

func IsSupportedFormat(format string) bool {
	return SupportedOutput[format]
}

func UpgradeGoPix(isUpgrade bool) {
	fmt.Println("\033[0;32m🔄 Starting GoPix upgrade...\033[0m")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("\033[0;32m❌ Failed to detect home directory:\033[0m", err)
		return
	}

	tmpDir := filepath.Join(homeDir, ".gopix_upgrade_tmp")
	repoURL := "https://github.com/MostafaSensei106/GoPix.git"

	// Check if repo already exists
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		fmt.Println("\033[0;32m📥 Cloning GoPix repository...\033[0m")
		cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("\033[0;32m❌ Failed to clone repository:\033[0m", err)
			return
		}
	} else {
		fmt.Println("\033[0;32m🔁 Repository exists. Checking for updates...\033[0m")
		cmd := exec.Command("git", "-C", tmpDir, "pull")
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("\033[0;32m❌ Failed to pull latest changes:\033[0m", err)
			return
		}
		if strings.Contains(output.String(), "Already up to date.") {
			fmt.Println("\033[0;32m✅ GoPix is already up to date.\033[0m")
			return
		} else {
			fmt.Println("\033[0;32m📦 Updates pulled successfully.\033[0m")
		}
	}

	fmt.Println("\033[0;32m🚀 Running install.sh ...\033[0m")
	installScript := filepath.Join(tmpDir, "install.sh")

	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("powershell.exe", installScript, "-f")
	} else {
		cmd = exec.Command("bash", installScript, "-f")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("\033[0;32m❌ Installation failed:\033[0m", err)
		return
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
