package upgrade

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/MostafaSensei106/GoPix/internal/platform"
)

const repoURL = "https://github.com/MostafaSensei106/GoPix.git"

// upgradeDirectory returns the path to the upgrade directory where
// the GoPix repository is cloned for upgrading. It is located in the
// user's home directory under ".gopix/upgrade". If the home directory
// cannot be determined, it logs an error and returns an empty string.

func upgradeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		color.Red("❌ Failed to detect home directory: %v", err)
		return ""
	}
	return filepath.Join(homeDir, ".gopix/upgrade")
}

// localRepositoryExists checks if the local GoPix repository exists
// in the specified upgrade directory. It returns true if the directory
// exists, otherwise it returns false.

func localRepositoryExists(upgradeDirectory string) bool {
	if _, err := os.Stat(upgradeDirectory); os.IsNotExist(err) {
		return false
	}

	return true

}

// cloneGitHubRepository clones the GoPix GitHub repository to the specified
// upgrade directory. It executes the command "git clone --depth=1 <repoURL>
// <upgradeDirectory>" and displays the output. If the command fails, it logs an
// error and returns.
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

// localRepositoryHash gets the hash of the local GoPix repository in the specified
// upgrade directory. If the command fails, it logs an error and returns an empty
// string.
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

// remoteRepositoryHash retrieves the latest commit hash of the remote GoPix
// repository from the specified repoURL. It executes the command
// "git ls-remote <repoURL> HEAD" and captures the output. If the command
// fails, it logs an error and returns an empty string. If successful, it
// returns the hash as a string.

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

// compareHashes checks if two hashes are equal.
//
// It takes two string arguments, a and b, and returns true if they are equal,
// false otherwise.
func compareHashes(a, b string) bool {
	return a == b
}

// pullLatestChanges pulls the latest changes from the remote GoPix repository to
// the specified upgrade directory. If the command fails, it logs an error and
// returns. If successful, it logs a success message.
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

// systemInfo returns a tuple containing the OS type and the CPU
// architecture type of the current system as strings.
//
// The first element of the tuple is the OS type, which can be one of
// "linux", "windows", or "darwin".
//
// The second element of the tuple is the CPU architecture type, which
// can be one of "amd64", "arm64", "386", "arm", "mips", "mipsle", or
// "mips64".
func systemInfo() (string, string) {
	return platform.OSType(), platform.ArchType()
}

// installUpdates installs the latest GoPix version by running the "make install"
// command inside the upgrade directory. It logs a message before and after
// installation, and errors if the command fails.
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

// UpgradeGoPix upgrades the GoPix application to the latest version.
// It checks if the local repository exists in the upgrade directory.
// If it does, it compares the local and remote repository hashes to
// determine if an update is required. If the hashes differ, it pulls
// the latest changes and installs the updates. If the local repository
// does not exist, it clones the repository and installs the updates.
// Informative messages are displayed throughout the process.

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
