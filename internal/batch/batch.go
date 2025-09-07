package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MostafaSensei106/GoPix/internal/config"
	"github.com/MostafaSensei106/GoPix/internal/logger"
	"github.com/MostafaSensei106/GoPix/internal/validator"
)

// BatchProcessor handles batch processing of folders and subfolders
type BatchProcessor struct {
	config *config.BatchConfig
}

// BatchResult contains information about a batch processing operation
type BatchResult struct {
	TotalFiles     int
	ProcessedFiles int
	FailedFiles    int
	SkippedFiles   int
	Directories    []string
	Errors         []error
}

// FileInfo contains information about a file to be processed
type FileInfo struct {
	Path      string
	RelPath   string // Relative path from input directory
	Dir       string // Directory containing the file
	Extension string
	Size      int64
}

// NewBatchProcessor creates a new BatchProcessor with the given configuration
func NewBatchProcessor(batchConfig *config.BatchConfig) *BatchProcessor {
	return &BatchProcessor{
		config: batchConfig,
	}
}

// CollectFilesRecursively collects all image files from the specified directory
// and its subdirectories based on the batch processing configuration
func (bp *BatchProcessor) CollectFilesRecursively(inputDir string, supportedExts []string) ([]FileInfo, error) {
	var files []FileInfo
	var mu sync.Mutex

	// Create a map for quick extension lookup
	extMap := make(map[string]bool)
	for _, ext := range supportedExts {
		extMap[strings.ToLower(ext)] = true
	}

	// Walk function to process each file/directory
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Log error but continue processing
			logger.Logger.Warnf("Error accessing path %s: %v", path, err)
			return nil
		}

		// Skip if it's a directory
		if info.IsDir() {
			return nil
		}

		// Check if we should follow symlinks
		if !bp.config.FollowSymlinks && info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Validate file path for security
		if err := validator.ValidateFilePath(path); err != nil {
			logger.Logger.Warnf("Skipping invalid path: %s", path)
			return nil
		}

		// Check file extension
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(info.Name()), "."))
		if !extMap[ext] {
			return nil
		}

		// Calculate relative path from input directory
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			logger.Logger.Warnf("Could not calculate relative path for %s: %v", path, err)
			relPath = filepath.Base(path)
		}

		// Check depth limit if configured
		if bp.config.MaxDepth > 0 {
			depth := strings.Count(relPath, string(filepath.Separator))
			if depth > bp.config.MaxDepth {
				return nil
			}
		}

		// Create file info
		fileInfo := FileInfo{
			Path:      path,
			RelPath:   relPath,
			Dir:       filepath.Dir(path),
			Extension: ext,
			Size:      info.Size(),
		}

		// Thread-safe append
		mu.Lock()
		files = append(files, fileInfo)
		mu.Unlock()

		return nil
	}

	// Walk the directory
	if err := filepath.Walk(inputDir, walkFunc); err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", inputDir, err)
	}

	return files, nil
}

// CollectFilesNonRecursive collects image files only from the specified directory
// without traversing subdirectories
func (bp *BatchProcessor) CollectFilesNonRecursive(inputDir string, supportedExts []string) ([]FileInfo, error) {
	var files []FileInfo

	// Create a map for quick extension lookup
	extMap := make(map[string]bool)
	for _, ext := range supportedExts {
		extMap[strings.ToLower(ext)] = true
	}

	// Read directory contents
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", inputDir, err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		// Check if we should follow symlinks
		if !bp.config.FollowSymlinks && entry.Type()&os.ModeSymlink != 0 {
			continue
		}

		path := filepath.Join(inputDir, entry.Name())

		// Validate file path for security
		if err := validator.ValidateFilePath(path); err != nil {
			logger.Logger.Warnf("Skipping invalid path: %s", path)
			continue
		}

		// Check file extension
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(entry.Name()), "."))
		if !extMap[ext] {
			continue
		}

		// Get file info
		info, err := entry.Info()
		if err != nil {
			logger.Logger.Warnf("Could not get file info for %s: %v", path, err)
			continue
		}

		// Create file info
		fileInfo := FileInfo{
			Path:      path,
			RelPath:   entry.Name(),
			Dir:       inputDir,
			Extension: ext,
			Size:      info.Size(),
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// CollectFiles collects image files based on the batch processing configuration
func (bp *BatchProcessor) CollectFiles(inputDir string, supportedExts []string) ([]FileInfo, error) {
	if bp.config.RecursiveSearch {
		return bp.CollectFilesRecursively(inputDir, supportedExts)
	}
	return bp.CollectFilesNonRecursive(inputDir, supportedExts)
}

// GetOutputPath calculates the output path for a file based on batch processing settings
func (bp *BatchProcessor) GetOutputPath(inputDir, filePath, targetFormat string) string {
	// Calculate relative path from input directory
	relPath, err := filepath.Rel(inputDir, filePath)
	if err != nil {
		// Fallback to just the filename
		relPath = filepath.Base(filePath)
	}

	// Remove original extension and add target extension
	basePath := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	newPath := basePath + "." + targetFormat

	// If custom output directory is specified, use it
	if bp.config.OutputDir != "" {
		return filepath.Join(bp.config.OutputDir, newPath)
	}

	// If preserving structure, maintain the relative path structure
	if bp.config.PreserveStructure {
		return filepath.Join(inputDir, newPath)
	}

	// Otherwise, put all files in the input directory
	return filepath.Join(inputDir, filepath.Base(newPath))
}

// GroupFilesByDirectory groups files by their containing directory
func (bp *BatchProcessor) GroupFilesByDirectory(files []FileInfo) map[string][]FileInfo {
	groups := make(map[string][]FileInfo)

	for _, file := range files {
		groups[file.Dir] = append(groups[file.Dir], file)
	}

	return groups
}

// GetDirectoryStats returns statistics about directories in the batch
func (bp *BatchProcessor) GetDirectoryStats(files []FileInfo) map[string]int {
	stats := make(map[string]int)

	for _, file := range files {
		stats[file.Dir]++
	}

	return stats
}

// ValidateBatchInput validates the input directory for batch processing
func (bp *BatchProcessor) ValidateBatchInput(inputDir string) error {
	// Check if input directory exists
	info, err := os.Stat(inputDir)
	if err != nil {
		return fmt.Errorf("input directory does not exist: %w", err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("input path is not a directory: %s", inputDir)
	}

	// Check if we have read permissions
	if info.Mode()&0400 == 0 {
		return fmt.Errorf("no read permission for directory: %s", inputDir)
	}

	return nil
}

// CreateOutputDirectory creates the output directory if it doesn't exist
func (bp *BatchProcessor) CreateOutputDirectory(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", dir, err)
	}
	return nil
}
