package platform

import "runtime"

// operatingSystem returns the current operating system.
//
// Possible values include:
//   - "linux"
//   - "windows"
//   - "darwin" (macOS)
func operatingSystem() string {
	return runtime.GOOS
}

// cpuArchitecture returns the current CPU architecture.
//
// Possible values include:
//   - "amd64" (x64)
//   - "arm64"
//   - "386" (x86)
//   - "arm"
//   - "mips"
//   - "mipsle"
//   - "mips64"
func cpuArchitecture() string {
	return runtime.GOARCH
}

func OSType() string {
	switch operatingSystem() {
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	case "Linux":
		return "macOS"
	default:
		return "Unknown"
	}
}

// ArchType returns the current CPU architecture type as a string.
//
// Possible values include:
//   - "amd64" (x64)
//
// If the architecture is not recognized, "Unknown" is returned.

func ArchType() string {
	switch cpuArchitecture() {
	case "amd64":
		return "amd64"
	default:
		return "Unknown"
	}
}
