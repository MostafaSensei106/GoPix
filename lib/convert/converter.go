package converter

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/chai2010/webp"
    "github.com/nfnt/resize"
    "golang.org/x/image/bmp"
    "golang.org/x/image/tiff"
)

type ConvertOptions struct {
    Quality      int
    MaxDimension int
    KeepOriginal bool
    DryRun       bool
    Backup       bool
}

type ConversionResult struct {
    OriginalPath string
    NewPath      string
    OriginalSize int64
    NewSize      int64
    Duration     time.Duration
    Error        error
}

type ImageConverter struct {
    options ConvertOptions
}

func NewImageConverter(options ConvertOptions) *ImageConverter {
    return &ImageConverter{options: options}
}

func (ic *ImageConverter) Convert(path, format string) *ConversionResult {
    start := time.Now()
    result := &ConversionResult{
        OriginalPath: path,
        Duration:     0,
    }

    defer func() {
        result.Duration = time.Since(start)
    }()

    // Get original file size
    if stat, err := os.Stat(path); err == nil {
        result.OriginalSize = stat.Size()
    }

    currentExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
    if currentExt == format || (currentExt == "jpg" && format == "jpeg") || (currentExt == "jpeg" && format == "jpg") {
        result.Error = fmt.Errorf("file already in target format")
        return result
    }

    newPath := strings.TrimSuffix(path, filepath.Ext(path)) + "." + format
    result.NewPath = newPath

    if ic.options.DryRun {
        return result
    }

    // Create backup if requested
    if ic.options.Backup {
        if err := ic.createBackup(path); err != nil {
            result.Error = fmt.Errorf("backup failed: %v", err)
            return result
        }
    }

    // Convert image
    if err := ic.convertImage(path, newPath, format); err != nil {
        result.Error = err
        return result
    }

    // Get new file size
    if stat, err := os.Stat(newPath); err == nil {
        result.NewSize = stat.Size()
    }

    // Remove original if not keeping
    if !ic.options.KeepOriginal {
        if err := os.Remove(path); err != nil {
            result.Error = fmt.Errorf("failed to remove original: %v", err)
            return result
        }
    }

    return result
}

func (ic *ImageConverter) convertImage(inputPath, outputPath, format string) error {
    file, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    // Decode image
    img, _, err := image.Decode(file)
    if err != nil {
        return fmt.Errorf("failed to decode image: %v", err)
    }

    // Resize if necessary
    if ic.options.MaxDimension > 0 {
        bounds := img.Bounds()
        if bounds.Dx() > ic.options.MaxDimension || bounds.Dy() > ic.options.MaxDimension {
            img = resize.Resize(uint(ic.options.MaxDimension), 0, img, resize.Lanczos3)
        }
    }

    // Create output file
    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer outFile.Close()

    // Encode based on format
    switch format {
    case "png":
        encoder := &png.Encoder{
            CompressionLevel: png.BestSpeed,
        }
        err = encoder.Encode(outFile, img)
    case "jpg", "jpeg":
        err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: ic.options.Quality})
    case "webp":
        err = webp.Encode(outFile, img, &webp.Options{
            Lossless: false,
            Quality:  float32(ic.options.Quality),
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
        return fmt.Errorf("failed to encode image: %v", err)
    }

    return nil
}

func (ic *ImageConverter) createBackup(path string) error {
    backupDir := filepath.Join(filepath.Dir(path), ".gopix_backup")
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return err
    }

    backupPath := filepath.Join(backupDir, filepath.Base(path))
    input, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    return os.WriteFile(backupPath, input, 0644)
}
