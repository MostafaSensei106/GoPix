package converter

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	// "golang.org/x/image/bmp"
)

// ConvertOptions contains the settings for the image conversion process.
type ConvertOptions struct {
	Quality      uint16
	MaxDimension uint16
	KeepOriginal bool
	DryRun       bool
	Backup       bool
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

	cache sync.Map
}

// cacheEntry stores metadata about converted images.
type cacheEntry struct {
	outputPath   string
	outputSize   int64
	lastModified time.Time
	configHash   string
}

// bufferPool manages reusable buffers to reduce GC pressure.
type bufferPool struct {
	ch chan []byte
}

func newBufferPool(size, capacity int) *bufferPool {
	bp := &bufferPool{
		ch: make(chan []byte, capacity),
	}
	// Pre-fill pool with buffers
	for i := 0; i < capacity; i++ {
		bp.ch <- make([]byte, size)
	}
	return bp
}

func (bp *bufferPool) get() []byte {
	select {
	case buf := <-bp.ch:
		return buf
	default:
		return make([]byte, 32*1024) // 32KB default
	}
}

func (bp *bufferPool) put(buf []byte) {
	if cap(buf) >= 32*1024 { // Only reuse reasonably sized buffers
		select {
		case bp.ch <- buf[:cap(buf)]:
		default:
			// Pool full, let GC handle it
		}
	}
}

// NewImageConverter returns a new ImageConverter instance.
func NewImageConverter(options ConvertOptions) *ImageConverter {
	return &ImageConverter{
		options: options,
		bufPool: newBufferPool(32*1024, 10), // 32KB buffers, pool of 10
		cache:   sync.Map{},
	}
}

// Convert converts the image at the given path to the given format.
func (ic *ImageConverter) Convert(path string, format string) *ConversionResult {
	start := time.Now()
	result := &ConversionResult{
		OriginalPath: path,
	}

	defer func() {
		result.Duration = time.Since(start)
	}()

	// Get original file info - use more efficient stat
	stat, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat file: %w", err)
		return result
	}
	result.OriginalSize = stat.Size()

	// Optimize string operations - avoid repeated allocations
	currentExt := getFileExtension(path)
	format = strings.ToLower(format)

	if isAlreadyInFormat(currentExt, format) {
		result.Error = fmt.Errorf("file already in target format")
		return result
	}

	// Pre-calculate new path using string builder for efficiency
	var pathBuilder strings.Builder
	basePath := strings.TrimSuffix(path, filepath.Ext(path))
	pathBuilder.Grow(len(basePath) + len(format) + 1)
	pathBuilder.WriteString(basePath)
	pathBuilder.WriteByte('.')
	pathBuilder.WriteString(format)
	result.NewPath = pathBuilder.String()

	// Check cache for existing conversion using sync.Map's Load method
	cacheKey := ic.getCacheKey(path, format)
	cached, exists := ic.cache.Load(cacheKey)

	if exists {
		// Type assertion for the cached value
		cachedEntry, ok := cached.(*cacheEntry)
		if !ok {
			// Handle unexpected type, remove invalid entry
			ic.cache.Delete(cacheKey)
		} else {
			if ic.isCacheValid(cachedEntry, stat.ModTime(), result.NewPath) {
				result.NewSize = cachedEntry.outputSize
				return result
			}
			// Remove invalid cache entry
			ic.cache.Delete(cacheKey)
		}
	}

	if newStat, err := os.Stat(result.NewPath); err == nil {
		result.NewSize = newStat.Size()

		// Store in cache using sync.Map's Store method
		ic.cache.Store(cacheKey, &cacheEntry{
			outputPath:   result.NewPath,
			outputSize:   result.NewSize,
			lastModified: time.Now(),
			configHash:   ic.getConfigHash(),
		})
	}

	if ic.options.DryRun {
		// For dry run, still check if conversion is needed using DecodeConfig
		needsResize, err := ic.checkIfResizeNeeded(path)
		if err != nil {
			result.Error = err
			return result
		}
		// Store in result for information (could extend ConversionResult if needed)
		_ = needsResize // Use the information as needed
		return result
	}

	// Create backup if requested
	if ic.options.Backup {
		if err := ic.createBackup(path); err != nil {
			result.Error = fmt.Errorf("backup failed: %w", err)
			return result
		}
	}

	// Convert image
	if err := ic.convertImageOptimized(path, result.NewPath, format); err != nil {
		result.Error = err
		return result
	}

	// Get new file size and update cache
	if newStat, err := os.Stat(result.NewPath); err == nil {
		result.NewSize = newStat.Size()

		// Store in cache using sync.Map's Store method
		ic.cache.Store(cacheKey, &cacheEntry{
			outputPath:   result.NewPath,
			outputSize:   result.NewSize,
			lastModified: stat.ModTime(),
			configHash:   ic.getConfigHash(),
		})
	}

	// Remove original if not keeping
	if !ic.options.KeepOriginal {
		if err := os.Remove(path); err != nil {
			result.Error = fmt.Errorf("failed to remove original: %w", err)
			return result
		}
	}

	return result
}

// getFileExtension efficiently extracts and normalizes file extension.
func getFileExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return strings.ToLower(ext[1:]) // Skip the dot
	}
	return ""
}

// isAlreadyInFormat checks if file is already in target format.
func isAlreadyInFormat(currentExt, targetFormat string) bool {
	if currentExt == targetFormat {
		return true
	}
	// Handle jpg/jpeg equivalence
	return (currentExt == "jpg" && targetFormat == "jpeg") ||
		(currentExt == "jpeg" && targetFormat == "jpg")
}

// checkIfResizeNeeded uses DecodeConfig to efficiently check dimensions without full decode.
func (ic *ImageConverter) checkIfResizeNeeded(inputPath string) (bool, error) {
	if ic.options.MaxDimension == 0 {
		return false, nil
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return false, fmt.Errorf("failed to open file for config: %w", err)
	}
	defer file.Close()

	// Use DecodeConfig for fast dimension checking without loading full image
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return false, fmt.Errorf("failed to decode config: %w", err)
	}

	maxDim := int(ic.options.MaxDimension)
	return config.Width > maxDim || config.Height > maxDim, nil
}

// getCacheKey generates a unique key for caching based on input path and target format.
func (ic *ImageConverter) getCacheKey(inputPath, format string) string {
	// Create hash from path, format, and relevant options
	hasher := md5.New()
	hasher.Write([]byte(inputPath))
	hasher.Write([]byte(format))
	hasher.Write([]byte(ic.getConfigHash()))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getConfigHash creates a hash of conversion settings for cache validation.
func (ic *ImageConverter) getConfigHash() string {
	var configBuilder strings.Builder
	configBuilder.WriteString(strconv.FormatUint(uint64(ic.options.Quality), 10))
	configBuilder.WriteByte('_')
	configBuilder.WriteString(strconv.FormatUint(uint64(ic.options.MaxDimension), 10))
	return configBuilder.String()
}

// isCacheValid checks if cached conversion is still valid.
func (ic *ImageConverter) isCacheValid(cached *cacheEntry, sourceModTime time.Time, expectedOutputPath string) bool {
	// Check if source file is newer than cache
	if sourceModTime.After(cached.lastModified) {
		return false
	}

	// Check if output file still exists
	if _, err := os.Stat(expectedOutputPath); err != nil {
		return false
	}

	// Check if conversion settings changed
	if cached.configHash != ic.getConfigHash() {
		return false
	}

	return true
}

// convertImageOptimized converts with DecodeConfig optimization and early dimension checking.
func (ic *ImageConverter) convertImageOptimized(inputPath, outputPath, format string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// First, check if we need to resize using DecodeConfig (fast)
	var needsResize bool
	var originalConfig image.Config

	if ic.options.MaxDimension > 0 {
		config, _, err := image.DecodeConfig(file)
		if err != nil {
			return fmt.Errorf("failed to decode config: %w", err)
		}
		originalConfig = config
		maxDim := int(ic.options.MaxDimension)
		needsResize = config.Width > maxDim || config.Height > maxDim

		// Reset file pointer for actual decode
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to seek file: %w", err)
		}
	}

	// Use buffered reader for better I/O performance
	bufferedReader := bufio.NewReaderSize(file, 64*1024)

	// Decode image with format hint for faster decoding
	img, imgFormat, err := image.Decode(bufferedReader)
	if err != nil {
		return fmt.Errorf("failed to decode image (%s): %w", imgFormat, err)
	}

	// Resize only if needed (we already know from DecodeConfig)
	if needsResize {
		// Calculate new dimensions maintaining aspect ratio
		var newWidth, newHeight uint
		maxDim := uint(ic.options.MaxDimension)

		if originalConfig.Width > originalConfig.Height {
			newWidth = maxDim
			newHeight = 0 // Let resize calculate height
		} else {
			newWidth = 0 // Let resize calculate width
			newHeight = maxDim
		}
		img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	}

	// Create output file with optimized flags
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		if cerr := outFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close output file: %w", cerr)
		}
	}()

	// Use buffered writer for better I/O performance
	bufferedWriter := bufio.NewWriterSize(outFile, 64*1024)
	defer func() {
		if ferr := bufferedWriter.Flush(); ferr != nil && err == nil {
			err = fmt.Errorf("failed to flush output: %w", ferr)
		}
	}()

	// Encode based on format with optimized settings
	switch strings.ToLower(format) {
	case "png":
		encoder := &png.Encoder{
			CompressionLevel: png.BestSpeed, // Faster compression
		}
		err = encoder.Encode(bufferedWriter, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(bufferedWriter, img, &jpeg.Options{
			Quality: int(ic.options.Quality),
		})
	case "webp":
		err = webp.Encode(bufferedWriter, img, &webp.Options{
			Lossless: false,
			Quality:  float32(ic.options.Quality),
		})
	// case "bmp":
	// 	err = bmp.Encode(bufferedWriter, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// createBackup creates a backup of the specified file in a directory named "backup"
// in the same directory as the original file. Optimized for performance.
func (ic *ImageConverter) createBackup(path string) error {
	dir := filepath.Dir(path)
	backupDir := filepath.Join(dir, "backup")

	// Create backup directory with proper permissions
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Build backup filename efficiently
	var filenameBuilder strings.Builder
	baseName := filepath.Base(path)
	filenameBuilder.Grow(len(baseName) + 4)
	filenameBuilder.WriteString(baseName)
	filenameBuilder.WriteString(".bak")

	backupPath := filepath.Join(backupDir, filenameBuilder.String())

	if err := ic.copyFileOptimized(path, backupPath); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// copyFileOptimized performs an optimized atomic file copy using buffer pool
func (ic *ImageConverter) copyFileOptimized(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer srcFile.Close()

	// Create temp file in same directory as destination for atomic rename
	tmpFile, err := os.CreateTemp(filepath.Dir(dst), ".tmp_"+filepath.Base(dst))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	tmpName := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpName) // Clean up on error
	}()

	// Get buffer from pool
	buf := ic.bufPool.get()
	defer ic.bufPool.put(buf)

	// Copy with buffered I/O using our buffer
	if _, err := io.CopyBuffer(tmpFile, srcFile, buf); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// Ensure data is written to disk
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// Close temp file before rename
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpName, dst); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
