package cmd

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MostafaSensei106/GoPix/internal/batch"
	"github.com/MostafaSensei106/GoPix/internal/config"
	"github.com/MostafaSensei106/GoPix/internal/converter"
	"github.com/MostafaSensei106/GoPix/internal/logger"
	"github.com/MostafaSensei106/GoPix/internal/progress"
	"github.com/MostafaSensei106/GoPix/internal/resume"
	"github.com/MostafaSensei106/GoPix/internal/stats"
	"github.com/MostafaSensei106/GoPix/internal/validator"
	"github.com/MostafaSensei106/GoPix/internal/worker"
)

var (
	Version   = "v1.5.2"
	BuildTime = time.Now().Format("2006-01-02 3:04:05pm")
	cfg       *config.Config

	// Command flags
	inputDir     string
	targetFormat string
	keepOriginal bool
	dryRun       bool
	verbose      bool
	workers      uint8
	quality      uint16
	maxDimension uint16
	backup       bool
	resumeFlag   bool
	rateLimit    float64
	logToFile    bool

	// Batch processing flags
	recursiveSearch   bool
	maxDepth          int
	preserveStructure bool
	outputDir         string
	groupByFolder     bool
	skipEmptyDirs     bool
	followSymlinks    bool
)

var rootCmd = &cobra.Command{
	Use:   "gopix",
	Short: "Advanced image converter with parallel processing write in Go",
	Long: `GoPix v1.5.3 - Professional Image Converter

Created by MostafaSensei106
GitHub: https://github.com/MostafaSensei106/GoPix`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Initialize logger
		logLevel := cfg.LogLevel
		if verbose {
			logLevel = "debug"
		}
		return logger.Initialize(logLevel, logToFile)
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		if resumeFlag {
			return handleResume()
		}

		// Apply config defaults if not set via flags
		if workers == 0 {
			workers = cfg.Workers
		}
		if quality == 0 {
			quality = cfg.Quality
		}
		if maxDimension == 0 {
			maxDimension = cfg.MaxDimension
		}
		if targetFormat == "" {
			targetFormat = cfg.DefaultFormat
		}

		// Validate inputs
		if err := validator.ValidateInputs(inputDir, targetFormat, cfg.Extentions); err != nil {
			return err
		}

		logger.Logger.Infof("Starting conversion: %s -> %s", inputDir, targetFormat)

		return runConversion()
	},
}

// runConversion handles the overall image conversion process. It collects all
// image files from the specified input directory, sets up the necessary
// resources such as the image converter and worker pool, and processes each
// file for conversion. The function supports resuming from previous sessions
// and updates the conversion state accordingly. It also tracks and reports
// progress and statistics throughout the process, and handles any errors that
// occur during conversion. On successful completion, it clears the resume
// state and logs the overall success of the conversion process.

func runConversion() error {
	// Create batch processor with configuration
	batchConfig := &config.BatchConfig{
		RecursiveSearch:   recursiveSearch,
		MaxDepth:          maxDepth,
		PreserveStructure: preserveStructure,
		OutputDir:         outputDir,
		GroupByFolder:     groupByFolder,
		SkipEmptyDirs:     skipEmptyDirs,
		FollowSymlinks:    followSymlinks,
	}

	// Override with config defaults if flags not set
	if !recursiveSearch && !preserveStructure && outputDir == "" && !groupByFolder && !skipEmptyDirs && !followSymlinks {
		batchConfig = &cfg.BatchProcessing
	}

	batchProcessor := batch.NewBatchProcessor(batchConfig)

	// Validate batch input
	if err := batchProcessor.ValidateBatchInput(inputDir); err != nil {
		return fmt.Errorf("batch input validation failed: %v", err)
	}

	// Collect all image files using batch processor
	fileInfos, err := batchProcessor.CollectFiles(inputDir, cfg.Extentions)
	if err != nil {
		return fmt.Errorf("failed to collect files: %v", err)
	}

	if len(fileInfos) == 0 {
		color.Yellow("⚠️  No supported image files found in: %s", inputDir)
		return nil
	}

	// Convert FileInfo to string paths for compatibility
	files := make([]string, len(fileInfos))
	for i, fileInfo := range fileInfos {
		files[i] = fileInfo.Path
	}

	color.Cyan("🔍 Found %d image files to process", len(files))

	// Show batch processing info
	if batchConfig.RecursiveSearch {
		color.Cyan("📁 Recursive search enabled (max depth: %d)", batchConfig.MaxDepth)
	}
	if batchConfig.PreserveStructure {
		color.Cyan("📂 Preserving directory structure")
	}
	if batchConfig.OutputDir != "" {
		color.Cyan("📤 Output directory: %s", batchConfig.OutputDir)
	}

	// Setup conversion state for resume capability
	sessionID := generateSessionID()
	conversionState := &resume.ConversionState{
		ProcessedFiles: []string{},
		StartTime:      time.Now(),
		InputDir:       inputDir,
		TargetFormat:   targetFormat,
		TotalFiles:     len(files),
		SessionID:      sessionID,
	}

	if cfg.ResumeEnabled {
		if err := resume.SaveState(conversionState); err != nil {
			logger.Logger.Warnf("Failed to save initial state: %v", err)
		}
	}

	// Setup converter
	converterOptions := converter.ConvertOptions{
		Quality:      quality,
		MaxDimension: maxDimension,
		KeepOriginal: keepOriginal,
		DryRun:       dryRun,
		Backup:       backup,
	}

	imageConverter := converter.NewImageConverter(converterOptions)

	// Setup worker pool
	pool := worker.NewWorkerPool(workers, imageConverter, rateLimit)

	// Setup progress tracking
	progressReporter := progress.NewProgressReporter(uint32(len(files)), "Converting images")
	statistics := stats.NewConversionStatistics()

	// Set batch processing flags in statistics
	statistics.BatchMode = true
	statistics.RecursiveSearch = batchConfig.RecursiveSearch
	statistics.PreserveStructure = batchConfig.PreserveStructure

	// Start processing
	pool.Start()
	defer pool.Stop()

	// Send jobs to worker pool
	go func() {
		for _, file := range files {
			// Calculate output path using batch processor
			outputPath := batchProcessor.GetOutputPath(inputDir, file, targetFormat)

			// Create output directory if needed
			if err := batchProcessor.CreateOutputDirectory(outputPath); err != nil {
				logger.Logger.Errorf("Failed to create output directory for %s: %v", file, err)
				continue
			}

			pool.AddJob(worker.Job{
				Path:       file,
				Format:     targetFormat,
				OutputPath: outputPath,
			})
		}
	}()

	// Process results
	processedCount := 0
	for processedCount < len(files) {
		select {
		case result := <-pool.Results():
			processedCount++

			// Update statistics
			statistics.AddResult(result)

			// Update progress
			if result.Error != nil {
				progressReporter.UpdateWithMessage(1, fmt.Sprintf("❌ %s", filepath.Base(result.OriginalPath)))
				logger.Logger.Errorf("Conversion failed: %s - %v", result.OriginalPath, result.Error)
			} else if result.NewSize == 0 {
				progressReporter.UpdateWithMessage(1, fmt.Sprintf("⏭️  %s", filepath.Base(result.OriginalPath)))
			} else {
				progressReporter.UpdateWithMessage(1, fmt.Sprintf("✅ %s", filepath.Base(result.OriginalPath)))
				logger.Logger.Infof("Converted: %s -> %s", result.OriginalPath, result.NewPath)
			}

			// Update resume state
			if cfg.ResumeEnabled {
				conversionState.ProcessedFiles = append(conversionState.ProcessedFiles, result.OriginalPath)
				if err := resume.SaveState(conversionState); err != nil {
					logger.Logger.Warnf("Failed to update state: %v", err)
				}
			}

		case <-time.After(30 * time.Second):
			logger.Logger.Warn("Processing timeout, continuing...")
		}
	}

	// Finish progress reporting
	progressReporter.Finish()

	// Print final statistics
	statistics.PrintReport()

	// Clear resume state on successful completion
	if cfg.ResumeEnabled {
		if err := resume.ClearState(); err != nil {
			logger.Logger.Warnf("Failed to clear state: %v", err)
		}
	}

	logger.Logger.Info("Conversion completed successfully")
	return nil
}

// handleResume attempts to load a saved conversion state and, if found, resumes the conversion from where it left off.
// It will print the saved state details and continue with the normal conversion process.
func handleResume() error {
	state, err := resume.LoadState()
	if err != nil {
		return fmt.Errorf("failed to load resume state: %v", err)
	}

	if state == nil {
		color.Yellow("⚠️  No previous conversion session found to resume")
		return nil
	}

	color.Cyan("🔄 Resuming conversion session from %v", state.StartTime.Format("2006-01-02 15:04:05"))
	color.Cyan("📁 Input directory: %s", state.InputDir)
	color.Cyan("🎯 Target format: %s", state.TargetFormat)
	color.Cyan("📊 Progress: %d/%d files processed", len(state.ProcessedFiles), state.TotalFiles)

	// Set variables from saved state
	inputDir = state.InputDir
	targetFormat = state.TargetFormat

	// Continue with normal conversion (it will skip already processed files)
	return runConversion()
}

// generateSessionID generates a random 8-byte session ID as a hexadecimal string.
func generateSessionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("❌ Error: %v", err)
		os.Exit(1)
	}
}

// init initializes the command-line interface by setting up the root command
// with various flags and configurations. It defines input/output flags such
// as the image folder path and target format, quality and processing flags
// like output quality and number of workers, and feature flags for backup
// and resumption of conversions. The function marks the path flag as required,
// sets the version template, and adds subcommands like the upgrade command.

func init() {
	// Input/Output flags
	rootCmd.Flags().StringVarP(&inputDir, "path", "p", "", "Path to the image folder (required)")
	rootCmd.Flags().StringVarP(&targetFormat, "to", "t", "", "Target format default: png (png, jpg, jpeg, webp)")
	rootCmd.Flags().BoolVar(&keepOriginal, "keep", false, "Keep original images after conversion")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without converting")

	// Quality and processing flags
	rootCmd.Flags().Uint16VarP(&quality, "quality", "q", 0, "Output quality (1-100, default 80)")
	rootCmd.Flags().Uint16Var(&maxDimension, "max-size", 0, "Maximum width/height in pixels default no limit")
	rootCmd.Flags().Uint8VarP(&workers, "workers", "w", 0, "Number of parallel workers Default: Max CPU Cores Available")
	rootCmd.Flags().Float64Var(&rateLimit, "rate-limit", 0, "Operations per second limit Default: No limit")

	// Feature flags
	rootCmd.Flags().BoolVar(&backup, "backup", false, "Create backup of original files")
	rootCmd.Flags().BoolVar(&resumeFlag, "resume", false, "Resume previous interrupted conversion")
	// rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.Flags().BoolVar(&logToFile, "log-file", false, "Save logs to file")

	// Batch processing flags
	rootCmd.Flags().BoolVar(&recursiveSearch, "recursive", true, "Search subdirectories recursively (default: true)")
	rootCmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Maximum directory depth to search (0 = unlimited)")
	rootCmd.Flags().BoolVar(&preserveStructure, "preserve-structure", true, "Preserve directory structure in output (default: true)")
	rootCmd.Flags().StringVar(&outputDir, "output-dir", "", "Custom output directory for batch processing")
	rootCmd.Flags().BoolVar(&groupByFolder, "group-by-folder", false, "Group results by source folder")
	rootCmd.Flags().BoolVar(&skipEmptyDirs, "skip-empty", true, "Skip directories with no images (default: true)")
	rootCmd.Flags().BoolVar(&followSymlinks, "follow-symlinks", false, "Follow symbolic links")

	// Mark required flags
	rootCmd.MarkFlagRequired("path")

	// Set version
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("GoPix {{.Version}}\n")

	rootCmd.AddCommand(upgradeCmd)
}
