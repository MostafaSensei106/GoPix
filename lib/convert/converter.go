package convert

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func ConvertImage(path, format string, keepOriginal, dryRun bool) error {
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
		fmt.Printf("[DRY-RUN] Would convert: %s -> %s\n", path, newPath)
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
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
	case "webp":
		err = webp.Encode(outFile, img, &webp.Options{Lossless: true, Quality: 100})
	case "avif":
		err = avif.Encode(outFile, img, &avif.Options{Quality: 100})
	case "tiff":
		err = tiff.Encode(outFile, img, &tiff.Options{Compression: 0, Predictor: true})
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
