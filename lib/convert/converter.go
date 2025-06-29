package convert

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	// "github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func ConvertImage(path, format string, keepOriginal, dryRun bool) error {
	currentExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if currentExt == format || (currentExt == "jpg" && format == "jpeg") || (currentExt == "jpeg" && format == "jpg") {
		return nil
	}

	newPath := strings.TrimSuffix(path, filepath.Ext(path)) + "." + format
	if dryRun {
		fmt.Printf("[DRY-RUN] Would convert: %s -> %s\n", path, newPath)
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var img image.Image
	switch currentExt {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "webp":
		img, err = webp.Decode(file)
	case "bmp":
		img, err = bmp.Decode(file)
	case "tiff":
		img, err = tiff.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}

	if err != nil {
		return err
	}

	outFile, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "png":
		encoder := &png.Encoder{
			CompressionLevel: png.BestSpeed,
		}
		err = encoder.Encode(outFile, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 85})
	case "webp":
		err = webp.Encode(outFile, img, &webp.Options{
			Lossless: false,
			Quality:  85,
		})
	case "tiff":
		err = tiff.Encode(outFile, img, &tiff.Options{
			Compression: tiff.LZW,
			Predictor:   false,
		})
	case "bmp":
		err = bmp.Encode(outFile, img)
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
