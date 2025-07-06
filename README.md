# GoPix v2.0 - Professional Image Converter

A powerful, feature-rich image conversion tool built in Go with advanced parallel processing and smart features.

## 🚀 Features

### Core Functionality
- **Multi-format support**: PNG, JPEG, WebP, TIFF, BMP
- **Parallel processing**: Utilize all CPU cores for maximum speed
- **Smart resume**: Continue interrupted conversions automatically
- **Progress tracking**: Real-time progress bars and ETA
- **Quality control**: Configurable compression and quality settings

### Advanced Features
- **Automatic backup**: Keep original files safe
- **Size optimization**: Intelligent resizing and compression
- **Batch operations**: Process multiple directories
- **Comprehensive statistics**: Detailed conversion reports
- **Configuration management**: Persistent settings
- **Rate limiting**: Control resource usage
- **Dry-run mode**: Preview changes before conversion

### Security & Reliability
- **Path validation**: Prevent directory traversal attacks
- **Permission checking**: Verify read/write access
- **Error recovery**: Robust error handling and retry logic
- **Disk space validation**: Check available space before conversion

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/MostafaSensei106/GoPix.git
cd GoPix

# Build the application
go build -o gopix main.go

# Install globally (optional)
go install
```

## 🎯 Quick Start

```bash
# Convert all images in a directory to PNG
gopix -p /path/to/images -t png

# Convert with custom quality and keep originals
gopix -p /path/to/images -t jpg -q 90 --keep

# Dry run to preview changes
gopix -p /path/to/images -t webp --dry-run

# Resume interrupted conversion
gopix --resume
```

## 📋 Usage Examples

### Basic Conversion
```bash
# Convert to WebP with high quality
gopix -p ./photos -t webp -q 95

# Convert with size limit
gopix -p ./photos -t jpg --max-size 1920

# Convert with backup
gopix -p ./photos -t png --backup
```

### Advanced Options
```bash
# Use 8 workers with rate limiting
gopix -p ./photos -t jpg -w 8 --rate-limit 5

# Verbose logging to file
gopix -p ./photos -t png -v --log-file

```

## ⚙️ Configuration

GoPix uses a YAML configuration file located at `~/.gopix/config.yaml`:

```yaml
default_format: "png"
quality: 85
workers: 8
max_dimension: 4096
log_level: "info"
auto_backup: false
resume_enabled: true
supported_extensions: ["jpg", "jpeg", "png", "webp", "bmp", "tiff"]
```

## 📊 Performance

GoPix v2.0 delivers exceptional performance:

- **Parallel processing**: 4-8x faster than sequential conversion
- **Memory efficient**: Optimized memory usage for large files
- **Smart caching**: Reduced I/O operations
- **Rate limiting**: Prevent system overload

### Benchmarks
- 1000 JPEG files (2MB avg): ~30 seconds on 8-core CPU
- 500 PNG files (5MB avg): ~45 seconds on 8-core CPU
- Memory usage: <100MB for most operations

## 🛡️ Security Features

- **Path validation**: Prevents directory traversal attacks
- **Permission checking**: Validates file system access
- **Safe defaults**: Secure configuration out of the box
- **Backup system**: Protects against data loss

## 🔧 Development

### Prerequisites
- Go 1.21 or later
- Git

### Building from Source
```bash
git clone https://github.com/MostafaSensei106/GoPix.git
cd GoPix
go mod download
go build -o gopix main.go
```

### Running Tests
```bash
go test ./...
```

## 📝 Command Reference

### Global Flags
- `-p, --path`: Input directory (required)
- `-t, --to`: Target format
- `-q, --quality`: Output quality (1-100)
- `-w, --workers`: Number of parallel workers
- `--keep`: Keep original files
- `--dry-run`: Preview changes only
- `-v, --verbose`: Enable verbose logging

### Commands
- `gopix`: Main conversion command
- `gopix config`: Manage configuration
- `gopix info`: Analyze directories
- `gopix clean`: Clean temporary files
- `gopix --resume`: Resume interrupted conversion

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.