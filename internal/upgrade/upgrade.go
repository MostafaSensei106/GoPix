package upgrade

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/mostafasensei106/gopix/internal/platform"
)

const repoURL = "https://github.com/MostafaSensei106/GoPix.git"

func upgradeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		color.Red("❌ Failed to detect home directory: %v", err)
		return ""
	}
	return filepath.Join(homeDir, ".gopix/upgrade")
}

func localRepositoryExists(upgradeDirectory string) bool {
	if _, err := os.Stat(upgradeDirectory); os.IsNotExist(err) {
		return false
	}

	return true

}

func cloneGitHubRepository(upgradeDirectory string) {
	color.Cyan("📥 Cloning GoPix GitHub Repository...")
	cmd := exec.Command("git", "clone", "--depth=1", repoURL, upgradeDirectory)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		color.Red("❌ Failed to clone repository: %v", err)
		return
	}
}

func localRepositoryHash(upgradeDirectory string) string {
	cmd := exec.Command("git", "-C", upgradeDirectory, "rev-parse", "HEAD")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		color.Red("❌ Failed to get local repository hash: %v", err)
		return ""
	}
	hash := strings.Fields(output.String())
	if len(hash) > 0 {
		color.Cyan("📦 Local repository hash: %s", hash[0])
		return hash[0]
	}
	return ""
}

func remoteRepositoryHash(repoURL string) string {
	cmd := exec.Command("git", "ls-remote", repoURL, "HEAD")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		color.Red("❌ Failed to get remote repository hash: %v", err)
		return ""
	}
	hash := strings.Fields(output.String())
	if len(hash) > 0 {
		color.Cyan("📦 Remote repository hash: %s", hash[0])
		return hash[0]
	}
	return ""
}

func compareHashes(a, b string) bool {
	return a == b
}

func pullLatestChanges(upgradeDirectory string) {
	color.Green("✅ New Updates available.")
	color.Cyan("🔁 Getting latest updates from GitHub...")
	cmd := exec.Command("git", "-C", upgradeDirectory, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	color.Green("📦 Updates pulled successfully.")
	if err := cmd.Run(); err != nil {
		color.Red("❌ Failed to pull latest changes: %v", err)
		return
	}
}

func systemInfo() (string, string) {
	return platform.OSType(), platform.ArchType()
}

func installUpdates() {
	color.Cyan("🚀 Installing Updates...")
	makeInstall := filepath.Join(upgradeDirectory(), "Makefile")

	var cmd *exec.Cmd
	osType, archType := systemInfo()

	cmd = exec.Command("make", "-f", makeInstall, "OS="+osType, "ARCH="+archType)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		color.Red("❌ Failed to install updates: %v", err)
		return
	}
}

func UpgradeGoPix() {
	color.Cyan("🔄 Starting GoPix upgrade...")

	path := upgradeDirectory()

	if localRepositoryExists(path) {
		if compareHashes(localRepositoryHash(path), remoteRepositoryHash(repoURL)) {
			color.Green("✅ You have the latest version of GoPix")

		} else {
			pullLatestChanges(path)
			installUpdates()
		}
	} else {
		cloneGitHubRepository(path)
		installUpdates()
	}
}
