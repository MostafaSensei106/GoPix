package converter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davidbyttow/govips/v2/vips"

	appErrors "github.com/MostafaSensei106/GoPix/internal/errors"
)

// ConvertOptions contains the settings for the image conversion process.
type ConvertOptions struct {
	Quality      uint16
	MaxDimension uint16
	KeepOriginal bool
	DryRun       bool
	Backup       bool
	Metadata     string
}

// ConversionResult holds the outcome of a single image conversion.
type ConversionResult struct {
	OriginalPath string
	NewPath      string
	OriginalSize int64
	NewSize      int64
	Duration     time.Duration
	Error        error
}

// ImageConverter is responsible for converting images.
type ImageConverter struct {
	options ConvertOptions
	bufPool *bufferPool
	cache   sync.Map
}

// cacheEntry stores metadata about converted images.
type cacheEntry struct {
	outputPath   string
	outputSize   int64
	lastModified time.Time
	configHash   string
}

// bufferPool manages reusable buffers to reduce GC pressure using sync.Pool
type bufferPool struct {
	pool *sync.Pool
}

func newBufferPool(size int) *bufferPool {
	return &bufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
}

func (bp *bufferPool) get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *bufferPool) put(buf []byte) {
	bp.pool.Put(&buf)
}

// NewImageConverter returns a new ImageConverter instance.
func NewImageConverter(options ConvertOptions) *ImageConverter {
	return &ImageConverter{
		options: options,
		bufPool: newBufferPool(32 * 1024), // 32KB buffers
		cache:   sync.Map{},
	}
}

// Convert converts the image at the given path to the given format.
func (ic *ImageConverter) Convert(path string, format string) *ConversionResult {
	return ic.ConvertWithOutputPath(path, format, "")
}

// ConvertWithOutputPath converts the image at the given path to the given format with a custom output path.
func (ic *ImageConverter) ConvertWithOutputPath(path string, format string, outputPath string) *ConversionResult {
	start := time.Now()
	result := &ConversionResult{
		OriginalPath: path,
	}

	defer func() {
		result.Duration = time.Since(start)
	}()

	stat, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat file: %w", err)
		return result
	}
	result.OriginalSize = stat.Size()

	currentExt := getFileExtension(path)
	format = strings.ToLower(format)

	if isAlreadyInFormat(currentExt, format) {
		result.Error = fmt.Errorf("file already in target format")
		return result
	}

	if outputPath != "" {
		result.NewPath = outputPath
	} else {
		basePath := strings.TrimSuffix(path, filepath.Ext(path))
		result.NewPath = basePath + "." + format
	}

	cacheKey := ic.getCacheKey(path, format)
	if cached, exists := ic.cache.Load(cacheKey); exists {
		cachedEntry, ok := cached.(*cacheEntry)
		if ok && ic.isCacheValid(cachedEntry, stat.ModTime(), result.NewPath) {
			result.NewSize = cachedEntry.outputSize
			return result
		}
		ic.cache.Delete(cacheKey)
	}

	if newStat, err := os.Stat(result.NewPath); err == nil {
		result.NewSize = newStat.Size()
		ic.cache.Store(cacheKey, &cacheEntry{
			outputPath:   result.NewPath,
			outputSize:   result.NewSize,
			lastModified: time.Now(),
			configHash:   ic.getConfigHash(),
		})
	}

	if ic.options.DryRun {
		return result
	}

	if ic.options.Backup {
		if err := ic.createBackup(path); err != nil {
			result.Error = fmt.Errorf("backup failed: %w", err)
			return result
		}
	}

	if err := ic.convertImage(path, result.NewPath, format); err != nil {
		result.Error = err
		return result
	}

	if newStat, err := os.Stat(result.NewPath); err == nil {
		result.NewSize = newStat.Size()
		ic.cache.Store(cacheKey, &cacheEntry{
			outputPath:   result.NewPath,
			outputSize:   result.NewSize,
			lastModified: stat.ModTime(),
			configHash:   ic.getConfigHash(),
		})
	}

	if !ic.options.KeepOriginal {
		if err := os.Remove(path); err != nil {
			result.Error = fmt.Errorf("failed to remove original: %w", err)
			return result
		}
	}

	return result
}

func (ic *ImageConverter) convertImage(inputPath, outputPath, format string) error {
	img, err := vips.NewImageFromFile(inputPath)
	if err != nil {
		return fmt.Errorf("%w: %w", appErrors.ErrCorruptedImage, err)
	}
	defer img.Close()

	// Resize if max dimension is set
	if ic.options.MaxDimension > 0 {
		maxDim := float64(ic.options.MaxDimension)
		scale := 1.0
		if img.Width() > img.Height() {
			if float64(img.Width()) > maxDim {
				scale = maxDim / float64(img.Width())
			}
		} else {
			if float64(img.Height()) > maxDim {
				scale = maxDim / float64(img.Height())
			}
		}
		if scale < 1.0 {
			if err := img.Resize(scale, vips.KernelLanczos3); err != nil {
				return fmt.Errorf("failed to resize image: %w", err)
			}
		}
	}

	// Get export parameters based on format
	params := ic.getExportParams(format)

	// Export the image to a byte buffer
	imgBytes, _, err := img.Export(params)
	if err != nil {
		return fmt.Errorf("failed to export image: %w", err)
	}

	// Write the buffer to the output file
	if err := os.WriteFile(outputPath, imgBytes, 0644); err != nil {
		return fmt.Errorf("failed to write image to file: %w", err)
	}

	return nil
}

func (ic *ImageConverter) getExportParams(format string) *vips.ExportParams {
	params := vips.NewDefaultExportParams()
	quality := int(ic.options.Quality)

	// Handle metadata stripping
	switch ic.options.Metadata {
	case "strip":
		params.StripMetadata = true
	case "keep":
		params.StripMetadata = false
	case "strip-location":
		// TODO: Implement selective stripping of location tags
		params.StripMetadata = false // For now, keep all metadata
	default:
		params.StripMetadata = false // Default to keeping metadata
	}

	switch format {
	case "png":
		params.Format = vips.ImageTypePNG
		params.Compression = 6 // Default compression
		params.Quality = quality
	case "jpg", "jpeg":
		params.Format = vips.ImageTypeJPEG
		params.Quality = quality
	case "webp":
		params.Format = vips.ImageTypeWEBP
		params.Quality = quality
	case "tiff":
		params.Format = vips.ImageTypeTIFF
		params.Quality = quality
	case "gif":
		params.Format = vips.ImageTypeGIF
		params.Quality = quality
	case "avif":
		params.Format = vips.ImageTypeAVIF
		params.Quality = quality
	case "heif":
		params.Format = vips.ImageTypeHEIF
		params.Quality = quality
	default:
		params.Format = vips.ImageTypeJPEG
	}
	return params
}

// getFileExtension efficiently extracts and normalizes file extension.
func getFileExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return strings.ToLower(ext[1:])
	}
	return ""
}

// isAlreadyInFormat checks if file is already in target format.
func isAlreadyInFormat(currentExt, targetFormat string) bool {
	if currentExt == targetFormat {
		return true
	}
	return (currentExt == "jpg" && targetFormat == "jpeg") ||
		(currentExt == "jpeg" && targetFormat == "jpg")
}

// getCacheKey generates a unique key for caching based on input path and target format.
func (ic *ImageConverter) getCacheKey(inputPath, format string) string {
	hasher := md5.New()
	hasher.Write([]byte(inputPath))
	hasher.Write([]byte(format))
	hasher.Write([]byte(ic.getConfigHash()))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getConfigHash creates a hash of conversion settings for cache validation.
func (ic *ImageConverter) getConfigHash() string {
	return strconv.FormatUint(uint64(ic.options.Quality), 10) + "_" + strconv.FormatUint(uint64(ic.options.MaxDimension), 10)
}

// isCacheValid checks if cached conversion is still valid.
func (ic *ImageConverter) isCacheValid(cached *cacheEntry, sourceModTime time.Time, expectedOutputPath string) bool {
	if sourceModTime.After(cached.lastModified) {
		return false
	}
	if _, err := os.Stat(expectedOutputPath); err != nil {
		return false
	}
	if cached.configHash != ic.getConfigHash() {
		return false
	}
	return true
}

// createBackup creates a backup of the specified file.
func (ic *ImageConverter) createBackup(path string) error {
	dir := filepath.Dir(path)
	backupDir := filepath.Join(dir, "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	baseName := filepath.Base(path)
	backupPath := filepath.Join(backupDir, baseName+".bak")
	return ic.copyFileOptimized(path, backupPath)
}

// copyFileOptimized performs an optimized atomic file copy.
func (ic *ImageConverter) copyFileOptimized(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer srcFile.Close()

	tmpFile, err := os.CreateTemp(filepath.Dir(dst), ".tmp_"+filepath.Base(dst))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	defer func() {
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}
	}()

	buf := ic.bufPool.get()
	defer ic.bufPool.put(buf)

	if _, err = io.CopyBuffer(tmpFile, srcFile, buf); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}
	if err = tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err = os.Rename(tmpFile.Name(), dst); err != nil {
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
