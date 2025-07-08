//go:build linux || darwin
// +build linux darwin

package validator

import (
	"syscall"
)

// HasSufficientSpace checks if the specified directory has enough available disk space.
// The function calls syscall.Statfs to obtain the amount of available disk space in bytes,
// and compares it to the specified requiredBytes. If the available space is greater than
// or equal to the required space, the function returns true; otherwise, it returns false.
func HasSufficientSpace(dir string, requiredBytes int64) bool {
    var stat syscall.Statfs_t
    if err := syscall.Statfs(dir, &stat); err != nil {
        return false
    }
    return int64(stat.Bavail)*int64(stat.Bsize) > requiredBytes
}