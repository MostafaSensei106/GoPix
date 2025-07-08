//go:build windows
// +build windows

package validator

import (
	"syscall"
	"unsafe"
)

// HasSufficientSpace checks if the specified directory has enough available space to store the specified number of bytes.
// Returns true if the directory has enough space, otherwise returns false.
func HasSufficientSpace(dir string, requiredBytes int64) bool {
	var freeBytes uint64
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceExW := kernel32.MustFindProc("GetDiskFreeSpaceExW")
	_, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dir))),
		0,
		0,
		uintptr(unsafe.Pointer(&freeBytes)),
	)
	if err != nil && err.Error() != "The operation completed successfully." {
		return false
	}

	return int64(freeBytes) > requiredBytes
}
