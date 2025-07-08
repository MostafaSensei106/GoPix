package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mostafasensei106/gopix/internal/converter"
)

type ConversionStatistics struct {
	TotalFiles       uint32
	ConvertedFiles   uint32
	SkippedFiles     uint32
	FailedFiles      uint32
	TotalSizeBefore  uint64
	TotalSizeAfter   uint64
	TotalDuration    time.Duration
	AverageDuration  time.Duration
	SpaceSaved       int
	CompressionRatio float64
	FailureReasons   map[string]uint32
}

// NewConversionStatistics creates a new ConversionStatistics instance, with the FailureReasons map initialized to hold conversion error reasons and counts.
func NewConversionStatistics() *ConversionStatistics {
	return &ConversionStatistics{
		FailureReasons: make(map[string]uint32),
	}
}

// AddResult increments the total number of files, total duration, and adds the size of the original and new files.
// If the result contains an error, it increments the failed files count and adds the error to the failure reasons map.
// If the result indicates that the file was skipped, it increments the skipped files count.
// Otherwise, it increments the converted files count and adds the original and new file sizes to the total sizes.
func (cs *ConversionStatistics) AddResult(result *converter.ConversionResult) {
	cs.TotalFiles++
	cs.TotalDuration += result.Duration

	if result.Error != nil {
		cs.FailedFiles++
		cs.FailureReasons[result.Error.Error()]++
		return
	}

	if result.OriginalPath == "" && result.NewSize == 0 {
		cs.SkippedFiles++
		return
	}

	cs.ConvertedFiles++
	cs.TotalSizeBefore += uint64(result.OriginalSize)
	cs.TotalSizeAfter += uint64(result.NewSize)

}

// Calculate computes the average duration, space saved, and compression ratio from the accumulated
// conversion results. It should be called after all results have been added to the ConversionStatistics
// instance.
func (cs *ConversionStatistics) Calculate() {
	if cs.TotalFiles > 0 {
		cs.AverageDuration = cs.TotalDuration / time.Duration(cs.TotalFiles)
	}

	cs.SpaceSaved = int(cs.TotalSizeBefore - cs.TotalSizeAfter)
	if cs.TotalSizeBefore > 0 {
		cs.CompressionRatio = float64(cs.TotalSizeAfter) / float64(cs.TotalSizeBefore)
	}
}

// PrintReport prints a summary of the conversion statistics to the console.
// It displays the total number of files processed, the number of converted,
// skipped, and failed files, the total conversion time, the average time per
// file, and the effective processing speed. It also displays the original and
// new total sizes of the files and the space saved (or increased) as a result
// of the conversion. Finally, it lists the failure reasons and the number of
// files that failed for each reason.
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
	color.Cyan("\n⏱️  Time Analysis")
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

	// Failure analysis
	if len(cs.FailureReasons) > 0 {
		color.Red("\n🔍 Failuer Analysis")
		for reason, count := range cs.FailureReasons {
			color.Red("  • %s: %d files", reason, count)
		}
	}
}

// formatBytes converts a size in bytes to a human-readable string using binary prefixes (e.g., KB, MB).
// It returns the size formatted with one decimal place and the appropriate unit, starting from bytes.
// For example, 1024 bytes is converted to "1.0 KB". This function supports units up to exabytes (EB).

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
