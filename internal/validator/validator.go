package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appErrors "github.com/MostafaSensei106/GoPix/internal/errors"
)

type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateInputs validates the input directory and target format.
//
// It checks if the input directory exists and has read permission, and if the target format is supported.
//
// If any of the checks fail, it returns a specific error type.
func ValidateInputs(inputDirectory, targetFormat string, supportedFormats []string) error {
	if inputDirectory == "" {
		return fmt.Errorf("%w: input directory is required", appErrors.ErrSourceNotFound)
	}

	if _, err := os.Stat(inputDirectory); os.IsNotExist(err) {
		return fmt.Errorf("%w: input directory %s does not exist", appErrors.ErrSourceNotFound, inputDirectory)
	}

	if !hasReadPermission(inputDirectory) {
		return fmt.Errorf("%w: input directory %s does not have read permission", appErrors.ErrPermissionDenied, inputDirectory)
	}

	if !isValidFormat(targetFormat, supportedFormats) {
		return fmt.Errorf("%w: target format %s is not supported", appErrors.ErrUnsupportedFormat, targetFormat)
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
	// For small lists, linear search is actually faster than map creation
	for _, supportedFormat := range supportedFormats {
		if format == supportedFormat {
			return true
		}
	}
	return false
}
