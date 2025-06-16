package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var SupportedOutput = map[string]bool{
	"png":  true,
	"jpg":  true,
	"jpeg": true,
	"webp": true,
}

func IsSupportedFormat(format string) bool {
	return SupportedOutput[format]
}

func upgrade() {
	fmt.Println("🚀 Starting GoPix upgrade...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Failed to detect home directory:", err)
		return
	}

	tmpDir := filepath.Join(homeDir, ".gopix_upgrade_tmp")
	repoURL := "https://github.com/MostafaSensei/GoPix.git" // ضع الرابط الحقيقي هنا

	os.RemoveAll(tmpDir)

	fmt.Println("⬇️ Cloning latest GoPix version...")
	cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed to clone repository:", err)
		return
	}

	fmt.Println("🔧 Running install.sh ...")
	installScript := filepath.Join(tmpDir, "install.sh")
	cmd = exec.Command("bash", installScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Installation failed:", err)
		return
	}

	fmt.Println("✅ GoPix upgraded successfully!")
}
