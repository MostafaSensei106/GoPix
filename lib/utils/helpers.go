package utils

var SupportedOutput = map[string]bool{
	"png":  true,
	"jpg":  true,
	"jpeg": true,
	"webp": true,
}

func IsSupportedFormat(format string) bool {
	return SupportedOutput[format]
}
