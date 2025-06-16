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
	"tiff": true,
	"avif": true,
}

func IsSupportedFormat(format string) bool {
	return SupportedOutput[format]
}

func UpgradeGoPix(isUpgrade bool) {
	fmt.Println("🔄 Starting GoPix upgrade...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Failed to detect home directory:", err)
		return
	}

	tmpDir := filepath.Join(homeDir, ".gopix_upgrade_tmp")
	repoURL := "https://github.com/MostafaSensei106/GoPix.git"

	// Check if repo already exists
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		fmt.Println("📥 Cloning GoPix repository...")
		cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("❌ Failed to clone repository:", err)
			return
		}
	} else {
		fmt.Println("🔁 Repository exists. Checking for updates...")
		cmd := exec.Command("git", "-C", tmpDir, "pull")
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("❌ Failed to pull latest changes:", err)
			return
		}
		if strings.Contains(output.String(), "Already up to date.") {
			fmt.Println("✅ GoPix is already up to date.")
			return
		} else {
			fmt.Println("📦 Updates pulled successfully.")
		}
	}

	fmt.Println("🚀 Running install.sh ...")
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
		fmt.Println("❌ Installation failed:", err)
		return
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
