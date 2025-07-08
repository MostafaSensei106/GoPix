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

// NewProgressReporter returns a new ProgressReporter with a progress bar that
// has the given total and description. The progress bar is configured with
// default options and a custom theme that displays a progress bar with a
// solid block (█) for the completed portion, a hollow block (░) for the
// remaining portion, a vertical bar (│) for the bar start and end, and a
// space ( ) for the bar padding. The progress bar also displays the count
// of items completed, the total count of items, the elapsed time, and an
// estimate of the remaining time.
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

// Update increments the progress bar and updates the current count.
// The progress bar is advanced by the specified increment.
//
// Parameters:
//   increment - The number of units to increment the progress by.
//
// Returns:
//   An error if the progress bar fails to update.

func (pr *ProgressReporter) Update(increment uint32) error {
	pr.current += increment
	return pr.bar.Add(int(increment))
}

// UpdateWithMessage increments the progress bar and updates the current count.
// The progress bar is advanced by the specified increment and the description
// is updated with the given message.
//
// Parameters:
//   increment - The number of units to increment the progress by.
//   message - The message to display in the progress bar.
//
// Returns:
//   An error if the progress bar fails to update.
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

// GetProgress returns the current progress and total number of units.
//
// Returns:
//   A tuple containing the current progress and total number of units.
func (pr *ProgressReporter) GetProgress() (uint32, uint32) {
	return pr.current, pr.total
}
