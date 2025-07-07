package platform

import "runtime"

func operatingSystem() string {
	return runtime.GOOS
}

func cpuArchitecture() string {
	return runtime.GOARCH
}

func OSType () string {
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

func ArchType () string {
	switch cpuArchitecture() {
	case "amd64":
		return "x64"
	case "arm64":
		return "arm64"
	case "386":
		return "x86"
	case "arm":
		return "arm"
	case "mips":
		return "mips"
	case "mipsle":
		return "mipsle"
	case "mips64":
		return "mips64"
	case "mips64le":
		return "mips64le"
	case "ppc64":
		return "ppc64"
	case "ppc64le":
		return "ppc64le"
	case "s390":
		return "s390"
	case "s390x":
		return "s390x"
	default:
		return "Unknown"
	}
}
