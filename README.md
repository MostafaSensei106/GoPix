<h1 align="center">GoPix</h1>
<p align="center">
  <img src="https://socialify.git.ci/MostafaSensei106/GoPix/image?custom_language=Go&font=KoHo&language=1&logo=https%3A%2F%2Favatars.githubusercontent.com%2Fu%2F138288138%3Fv%3D4&name=1&owner=1&pattern=Floating+Cogs&theme=Light" alt="GoPix Banner">
</p>

<p align="center">
  <strong>A high-performance, feature-rich image conversion CLI tool built in Go.</strong><br>
  Fast. Smart. Efficient. All from the terminal.
</p>

<p align="center">
  <a href="#about">About</a> •
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#technologies">Technologies</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

## About

Welcome to **GoPix** — a blazing-fast image conversion CLI tool built with Go and powered by `libvips` for extreme performance.
GoPix empowers developers, designers, and power users with efficient batch image conversions, intelligent file handling, and performance-oriented architecture. Whether you’re processing thousands of photos or optimizing a single folder, GoPix handles it with speed and precision.

---

## Features

### 🌟 Core Functionality
- **High-Performance Engine**: Powered by `libvips` for 4-8x faster conversions and lower memory usage.
- **Extensive Format Support**: `PNG`, `JPG`, `WEBP`, `TIFF`, `GIF`, `AVIF`, `HEIF`.
- **Parallel Processing**: Uses all CPU cores for maximum speed.
- **Real-time Progress Bar**: Track progress with count, ETA, and throughput.
- **Smart Resume**: Automatically resume interrupted conversion sessions.

### 🛠️ Advanced Capabilities
- **Metadata Control**: Keep or strip EXIF data to reduce file size or protect privacy.
- **Enhanced Batch Processing**: Process folders and subfolders with advanced options.
  - Recursive directory traversal with depth control.
  - Preserve or flatten directory structure.
  - Custom output directory support.
- **Quality and Sizing**:
  - Custom output quality (1-100).
  - Set max width/height for automatic resizing.
- **Dry-Run Mode**: Preview all changes without writing any files.
- **Backup Originals**: Automatically back up original files before converting.
- **Rate Limiting**: Limit operations per second to prevent system overload.
- **Detailed Reporting**: Get a full statistical report after each session.

### 🛡️ Security & Reliability
- **Path Validation**: Prevents directory traversal attacks.
- **Permission Checking**: Ensures files and directories are accessible.
- **Disk Space Validation**: Checks for sufficient disk space before starting.

---

## Installation

### ⚠️ IMPORTANT: New Dependency

GoPix v2.0 and later uses `libvips` for image processing. You **must** have `libvips` installed on your system for GoPix to work.

#### 🔧 Installing libvips

- **On macOS:**
  ```bash
  brew install vips
  ```
- **On Debian/Ubuntu:**
  ```bash
  sudo apt install libvips-dev
  ```
- **On Fedora:**
  ```bash
  sudo dnf install vips-devel
  ```
- **On Windows:**
  1. Download the latest `vips-dev-w64-all-x.y.z.zip` from the [libvips releases page](https://github.com/libvips/libvips/releases).
  2. Extract it to a location like `C:\vips-dev`.
  3. Add the `C:\vips-dev\bin` directory to your system's `PATH`.

---

## 📦 Easy Install (Linux / Windows)

Download the latest pre-built binary for your platform from the [Releases](https://github.com/MostafaSensei106/GoPix/releases) page.

### 🐧 Linux
Extract the archive
```bash
tar -xzf gopix-linux-amd64.vX.Y.Z.tar.gz
```

Move the binary to the local bin directory
```bash
sudo mv linux/amd64/gopix /usr/local/bin
```

Then you can test the tool by running:

```bash
gopix -v
```
---

### 🪟 Windows

1. Download `gopix-windows-amd64-vX.Y.Z.zip` from the [Releases](https://github.com/MostafaSensei106/GoPix/releases) page.
2. Extract the archive to a folder of your choice.
3. Add that folder to your **System PATH**.

Then you can test the tool by running:
```powershell
gopix -v
```
---

## 🏗️ Build from Source

Ensure you have `Go`, `git`, `make`, and `libvips` installed first.

```bash
git clone --depth 1 https://github.com/MostafaSensei106/GoPix.git
cd GoPix
make
```
This will compile and install GoPix locally.

---

### 🆙 Upgrading

To upgrade GoPix to the latest version, simply run:
```bash
gopix upgrade
```

---

## 🚀 Quick Start

```bash
# Convert all images in a directory to high-quality AVIF
gopix -p ./images -t avif -q 90
```

```bash
# Convert to JPEG, strip metadata, and keep originals
gopix -p ./images -t jpg --metadata strip --keep
```

---

## 📋 Usage Examples

### 🔁 Basic Conversion
```bash
# Convert to WebP with 95% quality
gopix -p ./photos -t webp -q 95
```

### ⚙️ Metadata Control
```bash
# Convert to PNG and remove all EXIF data
gopix -p ./photos -t png --metadata strip
```

### 🔄 Advanced Batch Processing
```bash
# Process all images recursively and save to a different directory
gopix -p ./source_images -t webp --output-dir ./converted_images --recursive
```

---

## Configuration

GoPix uses a YAML config file located at `~/.gopix/config.yaml` on Linux/macOS and `%USERPROFILE%\.gopix\config.yaml` on Windows.

### 🧾 Example Config:
```yaml
default_format: "avif"
quality: 90
workers: 8
max_dimension: 4096
log_level: "info"
metadata: "keep" # Can be: keep, strip
auto_backup: false
resume_enabled: true

# Batch processing configuration
batch_processing:
  recursive_search: true
  max_depth: 0
  preserve_structure: true
  output_dir: ""
  group_by_folder: false
  skip_empty_dirs: true
  follow_symlinks: false
```

All settings can be overridden using CLI flags.

---

## Technologies

| Technology            | Description                                                                 |
|------------------------|-----------------------------------------------------------------------------|
| 🧠 **Golang**            | [go.dev](https://go.dev) — The core language powering GoPix: fast and efficient |
| 🚀 **Govips**           | [davidbyttow/govips](https://github.com/davidbyttow/govips) — High-performance image processing via libvips |
| 🛠️ **Cobra (CLI)**       | [spf13/cobra](https://github.com/spf13/cobra) — CLI commands, flags, and UX |
| 🎨 **Fatih/color**       | [fatih/color](https://github.com/fatih/color) — Terminal text styling and coloring |
| 📉 **Progress bar**      | [schollz/progressbar](https://github.com/schollz/progressbar) — Beautiful terminal progress bar |
| 📦 **YAML config**       | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — Config file parser |
| 📜 **Logrus**            | [sirupsen/logrus](https://github.com/sirupsen/logrus) — Advanced logging framework |

---

## Contributing

Contributions are welcome! Please open an issue first to discuss any major changes.

---

## License

This project is licensed under the **GPL-3.0 License**.
See the [LICENSE](LICENSE) file for full details.
<p align="center">
  Made with ❤️ by <a href="https://github.com/MostafaSensei106">MostafaSensei106</a>
</p>
