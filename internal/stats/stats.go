package stats

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/MostafaSensei106/GoPix/internal/converter"
	appErrors "github.com/MostafaSensei106/GoPix/internal/errors"
)

type FailureAnalysis struct {
	Corrupted   uint32
	Permission  uint32
	Unsupported uint32
	Other       uint32
}

type ConversionStatistics struct {
	TotalFiles           uint32
	ConvertedFiles       uint32
	SkippedFiles         uint32
	FailedFiles          uint32
	TotalSizeBefore      uint64
	TotalSizeAfter       uint64
	TotalDuration        time.Duration
	AverageDuration      time.Duration
	SpaceSaved           int
	CompressionRatio     float64
	Failures             FailureAnalysis
	DirectoriesProcessed map[string]int
	BatchMode            bool
	RecursiveSearch      bool
	PreserveStructure    bool
}

func NewConversionStatistics() *ConversionStatistics {
	return &ConversionStatistics{
		DirectoriesProcessed: make(map[string]int, 50),
	}
}

func (cs *ConversionStatistics) AddResult(result *converter.ConversionResult) {
	cs.TotalFiles++
	cs.TotalDuration += result.Duration

	if result.Error != nil {
		cs.FailedFiles++
		switch {
		case errors.Is(result.Error, appErrors.ErrCorruptedImage):
			cs.Failures.Corrupted++
		case errors.Is(result.Error, appErrors.ErrPermissionDenied):
			cs.Failures.Permission++
		case errors.Is(result.Error, appErrors.ErrUnsupportedFormat):
			cs.Failures.Unsupported++
		default:
			cs.Failures.Other++
		}
		return
	}

	if result.OriginalPath == "" && result.NewSize == 0 {
		cs.SkippedFiles++
		return
	}

	cs.ConvertedFiles++
	cs.TotalSizeBefore += uint64(result.OriginalSize)
	cs.TotalSizeAfter += uint64(result.NewSize)

	if cs.BatchMode {
		dir := filepath.Dir(result.OriginalPath)
		cs.DirectoriesProcessed[dir]++
	}
}

func (cs *ConversionStatistics) Calculate() {
	if cs.TotalFiles > 0 {
		cs.AverageDuration = cs.TotalDuration / time.Duration(cs.TotalFiles)
	}

	if cs.TotalSizeBefore > 0 {
		cs.SpaceSaved = int(cs.TotalSizeBefore - cs.TotalSizeAfter)
		cs.CompressionRatio = float64(cs.TotalSizeAfter) / float64(cs.TotalSizeBefore)
	}
}

func (cs *ConversionStatistics) PrintReport() {
	cs.Calculate()

	color.Cyan("\n📊 Conversion Report")
	color.Cyan(strings.Repeat("=", 50))

	// File statistics
	color.Green("✅ Converted: %d", cs.ConvertedFiles)
	color.Yellow("⏭️ Skipped: %d", cs.SkippedFiles)
	color.Red("❌ Failed: %d", cs.FailedFiles)
	color.Cyan("📁 Total processed: %d", cs.TotalFiles)

	// Time statistics
	color.Cyan("\n⏱️ Time Analysis")
	color.Cyan(strings.Repeat("=", 50))
	color.White("🔄 Total conversion time (sum of all file durations): %v", cs.TotalDuration.Round(time.Millisecond))
	color.White("📊 Avg. time per file: ~%v (non-parallel)", cs.AverageDuration.Round(time.Millisecond))
	if cs.ConvertedFiles > 0 {
		rate := float64(cs.ConvertedFiles) / cs.TotalDuration.Seconds()
		color.White("⚡ Effective processing speed: %.1f files/sec", rate)
	}

	// Size statistics
	if cs.TotalSizeBefore > 0 {
		color.Cyan("\n💾 Size Analysis")
		color.Cyan(strings.Repeat("=", 50))
		color.White("🗂️ Original total size: %s", formatBytes(int64(cs.TotalSizeBefore)))
		color.White("🆕 New total size: %s", formatBytes(int64(cs.TotalSizeAfter)))

		if cs.SpaceSaved > 0 {
			color.Green("💰 Space saved: %s (%.1f%% reduction)",
				formatBytes(int64(cs.SpaceSaved)),
				(1-cs.CompressionRatio)*100)
		} else if cs.SpaceSaved < 0 {
			color.Red("📈 Size increased: %s (%.1f%% increase)",
				formatBytes(-int64(cs.SpaceSaved)),
				(cs.CompressionRatio-1)*100)
		}
	}

	// Batch processing information
	if cs.BatchMode {
		color.Cyan("\n📁 Batch Processing")
		color.Cyan(strings.Repeat("=", 50))
		if cs.RecursiveSearch {
			color.White("🔄 Recursive search: Enabled")
		} else {
			color.White("🔄 Recursive search: Disabled")
		}
		if cs.PreserveStructure {
			color.White("📂 Directory structure: Preserved")
		} else {
			color.White("📂 Directory structure: Flattened")
		}
		color.White("📊 Directories processed: %d", len(cs.DirectoriesProcessed))
	}

	// Failure analysis
	if cs.FailedFiles > 0 {
		color.Red("\n🔍 Failure Analysis")
		color.Red(strings.Repeat("=", 50))
		if cs.Failures.Corrupted > 0 {
			color.Red("  • Corrupted images: %d", cs.Failures.Corrupted)
		}
		if cs.Failures.Permission > 0 {
			color.Red("  • Permission errors: %d", cs.Failures.Permission)
		}
		if cs.Failures.Unsupported > 0 {
			color.Red("  • Unsupported formats: %d", cs.Failures.Unsupported)
		}
		if cs.Failures.Other > 0 {
			color.Red("  • Other errors: %d", cs.Failures.Other)
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := "KMGTPE"
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), units[exp])
}
