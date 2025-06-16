package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/mostafasensei106/gopix/lib/convert"
	"github.com/mostafasensei106/gopix/lib/utils"
	"github.com/spf13/cobra"
)

var (
	inputDir     string
	targetFormat string
	keepOriginal bool
	dryRun       bool
)

var converted, skipped, failed int
var mu sync.Mutex

var rootCmd = &cobra.Command{
	Use:   "imgconvert",
	Short: "Convert images in a directory to a specific format",
	Run: func(cmd *cobra.Command, args []string) {
		if inputDir == "" || !utils.IsSupportedFormat(targetFormat) {
			color.Red("❌ Invalid arguments. Use -h for help.")
			return
		}
		color.Cyan("🔄 Converting images in: %s", inputDir)

		_ = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(info.Name()), "."))
			if ext == targetFormat {
				mu.Lock()
				skipped++
				mu.Unlock()
				return nil
			}

			err = convert.ConvertImage(path, targetFormat, keepOriginal, dryRun)
			mu.Lock()
			if err != nil {
				color.Red("[FAIL] %s (%v)", path, err)
				failed++
			} else {
				color.Green("[OK] %s", path)
				converted++
			}
			mu.Unlock()
			return nil
		})

		color.Cyan("\n✅ Converted: %d | ⏭️ Skipped: %d | ❌ Failed: %d", converted, skipped, failed)
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&inputDir, "path", "p", "", "Path to the image folder")
	rootCmd.Flags().StringVarP(&targetFormat, "to", "t", "png", "Target format (png, jpg, jpeg, webp)")
	rootCmd.Flags().BoolVar(&keepOriginal, "keep", false, "Keep original images after conversion")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without converting")
	_ = rootCmd.Execute()
}
