package errors

import "errors"

var (
	ErrCorruptedImage    = errors.New("corrupted image")
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrSourceNotFound    = errors.New("source not found")
	ErrFatal           = errors.New("fatal error")
)
