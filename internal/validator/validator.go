package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}



func ValidateInputs(inputDirectory, targetFormat string, supportedFormatschan []string) error {

	if inputDirectory == "" {
		return &ValidationError{"inputDir", "input directory is required"}
	}

	if _, err := os.Stat(inputDirectory); os.IsNotExist(err) {
		return &ValidationError{
			Field:   "inputDirectory",
			Message: fmt.Sprintf("input directory %s does not exist", inputDirectory),
		}
	}

	if !hasReadPermission(inputDirectory) {
		return &ValidationError{
			Field:   "inputDirectory",
			Message: fmt.Sprintf("input directory %s does not have read permission", inputDirectory),
		}
	}

	if !isValidFormat(targetFormat, supportedFormatschan) {
		return &ValidationError{
			Field:   "targetFormat",
			Message: fmt.Sprintf("target format %s is not supported", targetFormat),
		}
	}
	return nil
}

// ValidateFilePath checks if the given path is valid and does not contain any path traversal.
// Returns an error if the path is invalid, otherwise returns nil.
func ValidateFilePath(path string) error {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid path (contains path traversal): %s", path)
	}
	return nil
}

// func HasSufficientSpace(dir string, requiredBytes int64) bool {
//     var stat syscall.Statfs_t
//     if err := syscall.Statfs(dir, &stat); err != nil {
//         return false
//     }
//     return int64(stat.Bavail)*int64(stat.Bsize) > requiredBytes
// }


// hasReadPermission checks if the specified path can be opened for reading.
// Returns true if the path can be opened, otherwise returns false.

func hasReadPermission(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// isValidFormat checks if the given format is present in the list of supported formats.
// Returns true if the format is supported, otherwise returns false.

func isValidFormat(format string, supportedFormats []string) bool {
	for _, supportedFormat := range supportedFormats {
		if format == supportedFormat {
			return true
		}
	}
	return false
}
