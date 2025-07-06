package stats

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mostafasensei106/gopix/internal/converter"
	"strings"
	"time"
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
	SpaceSaved       uint64
	CompressionRatio float64
	FailureReasons   map[string]uint32
}

func NewConversionStatistics() *ConversionStatistics {
	return &ConversionStatistics{
		FailureReasons: make(map[string]uint32),
	}
}

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

func (cs *ConversionStatistics) Calculate() {
	if cs.TotalFiles > 0 {
		cs.AverageDuration = cs.TotalDuration / time.Duration(cs.TotalFiles)
	}

	cs.SpaceSaved = cs.TotalSizeBefore - cs.TotalSizeAfter
	if cs.TotalSizeBefore > 0 {
		cs.CompressionRatio = float64(cs.TotalSizeAfter) / float64(cs.TotalSizeBefore)
	}
}

func (cs *ConversionStatistics) PrintReport() {
	cs.Calculate()

	color.Cyan("\n📊 CONVERSION REPORT")
	color.Cyan(strings.Repeat("=", 50))

	// File statistics
	color.Green("✅ Converted: %d", cs.ConvertedFiles)
	color.Yellow("⏭️  Skipped: %d", cs.SkippedFiles)
	color.Red("❌ Failed: %d", cs.FailedFiles)
	color.Cyan("📁 Total processed: %d", cs.TotalFiles)

	// Size statistics
	if cs.TotalSizeBefore > 0 {
		color.Cyan("\n💾 SIZE ANALYSIS")
		color.White("Original size: %s", formatBytes(int64(cs.TotalSizeBefore)))
		color.White("New size: %s", formatBytes(int64(cs.TotalSizeAfter)))

		if cs.SpaceSaved > 0 {
			color.Green("💰 Space saved: %s (%.1f%% reduction)",
				formatBytes(int64(cs.SpaceSaved)),
				(1-cs.CompressionRatio)*100)
		} else if int64(cs.SpaceSaved) < 0 {
			color.Red("📈 Size increased: %s", formatBytes(-int64(cs.SpaceSaved)))
		}
	}

	// Time statistics
	color.Cyan("\n⏱️  PERFORMANCE")
	color.White("Total time: %v", cs.TotalDuration.Round(time.Millisecond))
	color.White("Average per file: %v", cs.AverageDuration.Round(time.Millisecond))
	if cs.ConvertedFiles > 0 {
		rate := float64(cs.ConvertedFiles) / cs.TotalDuration.Seconds()
		color.White("Processing rate: %.1f files/sec", rate)
	}

	// Failure analysis
	if len(cs.FailureReasons) > 0 {
		color.Red("\n🔍 FAILURE ANALYSIS")
		for reason, count := range cs.FailureReasons {
			color.Red("  • %s: %d files", reason, count)
		}
	}
}

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
