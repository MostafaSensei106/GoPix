package progress

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ProgressReporter struct {
	bar       *progressbar.ProgressBar
	startTime time.Time
	total     uint32
	current   uint32
}

func NewProgressReporter(total uint32, description string) *ProgressReporter {
	bar := progressbar.NewOptions(
		int(total),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(
			progressbar.Theme{
				Saucer:        "█",
				SaucerHead:    "█",
				SaucerPadding: "░",
				BarStart:      "│",
				BarEnd:        "│",
			}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetRenderBlankState(true),
	)
	return &ProgressReporter{
		bar:       bar,
		startTime: time.Now(),
		total:     total,
		current:   0,
	}
}

func (pr *ProgressReporter) Update(increment uint32) error {
	pr.current += increment
	return pr.bar.Add(int(increment))
}

func (pr *ProgressReporter) UpdateWithMessage(increment uint32, message string) error {
	pr.current += increment
	pr.bar.Describe(message)
	return pr.bar.Add(int(increment))
}

// Finish marks the progress bar as finished and prints the total elapsed time.
func (pr *ProgressReporter) Finish() {
	pr.bar.Finish()
	elapsed := time.Since(pr.startTime)
	fmt.Printf("\n⏱️  Wall-clock time (actual time from start to finish): %v\n", elapsed.Round(time.Millisecond))
}

func (pr *ProgressReporter) GetProgress() (uint32, uint32) {
	return pr.current, pr.total
}
