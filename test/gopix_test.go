
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MostafaSensei106/GoPix/internal/batch"
	"github.com/MostafaSensei106/GoPix/internal/config"
	"github.com/MostafaSensei106/GoPix/internal/converter"
	"github.com/MostafaSensei106/GoPix/internal/logger"
	"github.com/MostafaSensei106/GoPix/internal/platform"
	"github.com/MostafaSensei106/GoPix/internal/progress"
	"github.com/MostafaSensei106/GoPix/internal/resume"
	"github.com/MostafaSensei106/GoPix/internal/stats"
	"github.com/MostafaSensei106/GoPix/internal/validator"
	"github.com/MostafaSensei106/GoPix/internal/worker"
	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		err := logger.Initialize("info", false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logger.Logger.GetLevel() != logrus.InfoLevel {
			t.Errorf("expected log level info, got %v", logger.Logger.GetLevel())
		}
	})
}

func TestPlatform(t *testing.T) {
	t.Run("OSType", func(t *testing.T) {
		expected := "Unknown"
		switch runtime.GOOS {
		case "linux":
			expected = "Linux"
		case "windows":
			expected = "Windows"
		case "darwin":
			expected = "macOS"
		}
		if platform.OSType() != expected {
			t.Errorf("expected %s, got %s", expected, platform.OSType())
		}
	})

	t.Run("ArchType", func(t *testing.T) {
		expected := "Unknown"
		if runtime.GOARCH == "amd64" {
			expected = "amd64"
		}
		if platform.ArchType() != expected {
			t.Errorf("expected %s, got %s", expected, platform.ArchType())
		}
	})
}

func TestProgress(t *testing.T) {
	t.Run("NewProgressReporter", func(t *testing.T) {
		pr := progress.NewProgressReporter(100, "testing")
		if pr == nil {
			t.Fatal("expected progress reporter, got nil")
		}
		current, total := pr.GetProgress()
		if current != 0 || total != 100 {
			t.Errorf("expected 0/100, got %d/%d", current, total)
		}
	})
}

func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		cfg := config.DefaultConfig()
		if cfg.DefaultFormat != "png" {
			t.Errorf("expected png, got %s", cfg.DefaultFormat)
		}
	})
}

func TestValidator(t *testing.T) {
	t.Run("ValidateFilePath", func(t *testing.T) {
		if err := validator.ValidateFilePath("/safe/path"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := validator.ValidateFilePath("../unsafe/path"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("HasSufficientSpace", func(t *testing.T) {
		if !validator.HasSufficientSpace(".", 1) {
			t.Error("expected true, got false")
		}
	})
}

func TestResume(t *testing.T) {
	t.Run("SaveAndLoadState", func(t *testing.T) {
		state := &resume.ConversionState{
			SessionID: "test",
		}
		err := resume.SaveState(state)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		loaded, err := resume.LoadState()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if loaded.SessionID != "test" {
			t.Errorf("expected test, got %s", loaded.SessionID)
		}
		err = resume.ClearState()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestStats(t *testing.T) {
	t.Run("NewConversionStatistics", func(t *testing.T) {
		st := stats.NewConversionStatistics()
		if st == nil {
			t.Fatal("expected stats, got nil")
		}
	})
}

func TestBatch(t *testing.T) {
	t.Run("NewBatchProcessor", func(t *testing.T) {
		bp := batch.NewBatchProcessor(&config.BatchConfig{})
		if bp == nil {
			t.Fatal("expected batch processor, got nil")
		}
	})
}

func TestConverter(t *testing.T) {
	t.Run("NewImageConverter", func(t *testing.T) {
		ic := converter.NewImageConverter(converter.ConvertOptions{})
		if ic == nil {
			t.Fatal("expected image converter, got nil")
		}
	})
}

func TestWorker(t *testing.T) {
	t.Run("NewWorkerPool", func(t *testing.T) {
		wp := worker.NewWorkerPool(1, nil, 0)
		if wp == nil {
			t.Fatal("expected worker pool, got nil")
		}
	})
}

func TestAll(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "gopix_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test Config
	t.Run("Config", func(t *testing.T) {
		cfg := config.DefaultConfig()
		if cfg.Quality != 80 {
			t.Errorf("Expected default quality 80, got %d", cfg.Quality)
		}
	})

	// Test Logger
	t.Run("Logger", func(t *testing.T) {
		if err := logger.Initialize("debug", false); err != nil {
			t.Fatalf("Failed to initialize logger: %v", err)
		}
		if logger.Logger.GetLevel() != logrus.DebugLevel {
			t.Errorf("Expected debug log level, got %s", logger.Logger.GetLevel())
		}
	})

	// Test Validator
	t.Run("Validator", func(t *testing.T) {
		supportedFormats := []string{"png", "jpg"}
		if err := validator.ValidateInputs(tmpDir, "png", supportedFormats); err != nil {
			t.Errorf("Validation failed for valid inputs: %v", err)
		}
		if err := validator.ValidateInputs("nonexistent", "png", supportedFormats); err == nil {
			t.Error("Expected error for nonexistent directory")
		}
		if err := validator.ValidateInputs(tmpDir, "gif", supportedFormats); err == nil {
			t.Error("Expected error for unsupported format")
		}
	})

	// Test Batch Processor
	t.Run("BatchProcessor", func(t *testing.T) {
		// Create some test files
		testFile1 := filepath.Join(tmpDir, "test1.png")
		if err := os.WriteFile(testFile1, []byte("fake png"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		subDir := filepath.Join(tmpDir, "sub")
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create sub dir: %v", err)
		}
		testFile2 := filepath.Join(subDir, "test2.jpg")
		if err := os.WriteFile(testFile2, []byte("fake jpg"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		batchConfig := &config.BatchConfig{RecursiveSearch: true, PreserveStructure: true}
		bp := batch.NewBatchProcessor(batchConfig)
		files, err := bp.CollectFiles(tmpDir, []string{"png", "jpg"})
		if err != nil {
			t.Fatalf("Failed to collect files: %v", err)
		}
		if len(files) != 2 {
			t.Errorf("Expected 2 files, got %d", len(files))
		}

		outputPath := bp.GetOutputPath(tmpDir, testFile2, "webp")
		expectedPath := filepath.Join(tmpDir, "sub", "test2.webp")
		// Normalize paths for comparison
		if filepath.ToSlash(outputPath) != filepath.ToSlash(expectedPath) {
			t.Errorf("Expected output path %s, got %s", expectedPath, outputPath)
		}
	})

	// Test Worker Pool & Converter
	t.Run("WorkerPool and Converter", func(t *testing.T) {
		// This is a more complex integration test.
		// We'll test with a simple dry run.
		opts := converter.ConvertOptions{DryRun: true}
		ic := converter.NewImageConverter(opts)
		pool := worker.NewWorkerPool(1, ic, 0)
		pool.Start()

		testFile := filepath.Join(tmpDir, "job.png")
		if err := os.WriteFile(testFile, []byte("fake png"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			result := <-pool.Results()
			if result.Error != nil && !strings.Contains(result.Error.Error(), "corrupted") {
				t.Errorf("Conversion failed unexpectedly: %v", result.Error)
			}
			if result.OriginalPath != testFile {
				t.Errorf("Expected original path %s, got %s", testFile, result.OriginalPath)
			}
		}()

		pool.AddJob(worker.Job{Path: testFile, Format: "jpg"})
		pool.Stop()
		wg.Wait()
	})

	// Test Stats
	t.Run("Statistics", func(t *testing.T) {
		stats := stats.NewConversionStatistics()
		result := &converter.ConversionResult{
			OriginalSize: 1000,
			NewSize:      500,
			Duration:     100 * time.Millisecond,
		}
		stats.AddResult(result)
		stats.Calculate()

		if stats.ConvertedFiles != 1 {
			t.Errorf("Expected 1 converted file, got %d", stats.ConvertedFiles)
		}
		if stats.TotalSizeBefore != 1000 {
			t.Errorf("Expected total size before 1000, got %d", stats.TotalSizeBefore)
		}
		if stats.SpaceSaved != 500 {
			t.Errorf("Expected space saved 500, got %d", stats.SpaceSaved)
		}
	})

	// Test Resume
	t.Run("Resume", func(t *testing.T) {
		stateFile := filepath.Join(os.TempDir(), "gopix_resume_test.json")
		// Mock getStateDir to use a temp file
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", os.TempDir())
		defer os.Setenv("HOME", oldHome)
		os.Remove(stateFile) // Clean up previous runs

		state := &resume.ConversionState{
			InputDir:     tmpDir,
			TargetFormat: "jpg",
			TotalFiles:   10,
		}
		if err := resume.SaveState(state); err != nil {
			t.Fatalf("Failed to save state: %v", err)
		}

		loadedState, err := resume.LoadState()
		if err != nil {
			t.Fatalf("Failed to load state: %v", err)
		}
		if loadedState.InputDir != tmpDir {
			t.Errorf("Expected input dir %s, got %s", tmpDir, loadedState.InputDir)
		}

		if err := resume.ClearState(); err != nil {
			t.Fatalf("Failed to clear state: %v", err)
		}

		if _, err := os.Stat(stateFile); !os.IsNotExist(err) {
			t.Error("Expected state file to be cleared")
		}
	})
}
func TestMain(m *testing.M) {
	// Setup anything needed for all tests
	fmt.Println("Starting GoPix tests...")
	exitCode := m.Run()
	// Teardown
	fmt.Println("Finished GoPix tests.")
	os.Exit(exitCode)
}

