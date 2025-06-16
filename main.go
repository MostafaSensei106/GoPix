package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var (
	inputDir     string
	targetFormat string
	keepOriginal bool

dryRun       bool
)

var supportedOutput = map[string]bool{
	"png":  true,
	"jpg":  true,
	"jpeg": true,
	"webp": true,
}

var converted, skipped, failed int
var mu sync.Mutex

func convertImage(path, format string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	newPath := strings.TrimSuffix(path, filepath.Ext(path)) + "." + format
	if dryRun {
		color.Yellow("[DRY-RUN] Would convert: %s -> %s", path, newPath)
		return nil
	}

	outFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "png":
		err = png.Encode(outFile, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	case "webp":
		err = webp.Encode(outFile, img, &webp.Options{Lossless: true})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return err
	}

	if !keepOriginal {
		_ = os.Remove(path)
	}

	return nil
}

func processDir(path string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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

		err = convertImage(path, targetFormat)
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
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "imgconvert",
		Short: "Convert images in a directory to a specific format",
		Run: func(cmd *cobra.Command, args []string) {
			if inputDir == "" || !supportedOutput[targetFormat] {
				color.Red("Invalid arguments. Use -h for help.")
				return
			}
			color.Cyan("🔄 Converting images in: %s", inputDir)
			processDir(inputDir)
			color.Cyan("\n✅ Converted: %d | ⏭️ Skipped: %d | ❌ Failed: %d", converted, skipped, failed)
		},
	}

	rootCmd.Flags().StringVarP(&inputDir, "path", "p", "", "Path to the image folder")
	rootCmd.Flags().StringVarP(&targetFormat, "to", "t", "png", "Target format (png, jpg, jpeg, webp)")
	rootCmd.Flags().BoolVar(&keepOriginal, "keep", false, "Keep original images after conversion")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without converting")

	_ = rootCmd.Execute()
}
