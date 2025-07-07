//go:build linux || darwin
// +build linux darwin

package validator

import (
	"syscall"
)

func HasSufficientSpace(dir string, requiredBytes int64) bool {
    var stat syscall.Statfs_t
    if err := syscall.Statfs(dir, &stat); err != nil {
        return false
    }
    return int64(stat.Bavail)*int64(stat.Bsize) > requiredBytes
}