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
  <a href="#quick-start">Quick Start</a> •
  <a href="#usage-examples">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#technologies">Technologies</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

## 📌 About

Welcome to **GoPix** — a blazing-fast image conversion CLI tool built with Go.  
GoPix empowers developers, designers, and power users with efficient batch image conversions, intelligent file handling, and performance-oriented architecture. Whether you’re processing thousands of photos or optimizing a single folder, GoPix handles it with speed and precision.

---

## ✨ Features

### 🌟 Core Functionality
- Multi-format support: PNG, JPG, WebP, BMP, TIFF, and more
- Parallel processing: Uses all CPU cores for maximum speed
- Real-time progress bar with ETA
- Smart resume for interrupted conversions
- Custom quality and compression settings

### 🛠️ Advanced Capabilities
- Batch processing for folders and subfolders
- Size and resolution limits
- Configuration profiles with YAML support
- Dry-run mode to preview changes
- Automatic backup of originals
- Rate limiting to prevent system overload
- Detailed post-process stats and reporting

### 🛡️ Security & Reliability
- Path validation to prevent directory traversal
- Safe defaults and permission checking
- Disk space validation before starting jobs
- Robust error handling and auto-retry mechanism

---

## 🛠️ Installation

### 📦 Easy Install (Linux)

Download the latest binary `gopix` from the [Releases](https://github.com/MostafaSensei106/GoPix/releases) page and move it to:

```bash
/home/$USER/.local/bin/
```

### 🧪 Installation Script Makefile (All Platforms: Linux, Windows, macOS)

> [!NOTE]
> The script may modify environment-specific paths depending on the system.

> [!IMPORTANT]
> ```MacOS``` should (untested) 😃
> ```Windows``` support will be added later. 


```bash
git clone --depth 1 https://github.com/MostafaSensei106/GoPix.git
cd GoPix
make install
```

This will install GoPix on your system and place the binary in the appropriate executable path.

---

## ⚡ Quick Start

```bash
# Convert all images in a directory to PNG
gopix -p ./images -t png
```

```bash
# Convert to JPEG with 90% quality and keep originals
gopix -p ./images -t jpg -q 90 --keep
```

```bash
# Preview changes without applying them
gopix -p ./images -t webp --dry-run
```

```bash
# Resume a previously interrupted conversion job
gopix --resume
```

---

## 📋 Usage Examples

### 🔁 Basic Conversion
```bash
gopix -p ./photos -t webp -q 95
```

### 💾 With Backup
```bash
gopix -p ./photos -t png --backup
```

### ⚙️ Advanced Usage
```bash
gopix -p ./photos -t jpg -w 8 --rate-limit 5
gopix -p ./photos -t png -v --log-file
```

---

## ⚙️ Configuration

GoPix uses a YAML config file located at:

```bash
# on Linux 
~/Home/$USER/.gopix/config.yaml
```

### 🧾 Example Config:
```yaml
default_format: "png"
quality: 85
workers: 8
max_dimension: 4096
log_level: "info"
auto_backup: false
resume_enabled: true
# supported_extensions: ["jpg", "jpeg", "png", "webp"] # Do not add any formats here, 
```

All settings can be overridden using CLI flags.

---

## 🧪 Development

### 📚 Prerequisites
- Go 1.21 or later
- Git

### 🏗️ Build from Source

To compile the GoPix source code, follow these steps:

1. **Clone the Repository**:  
   Ensure you have `git` installed and clone the GoPix repository:
   ```bash
   git clone --depth https://github.com/MostafaSensei106/GoPix.git
   cd GoPix
   ```

2. **Download Dependencies**:  
   Use the `go mod` command to download all necessary dependencies:
   ```bash
   go mod download
   ```

3. **Build the Application**:  
   Compile the GoPix application using the `go build` command. This will generate an executable named `gopix`:
   ```bash
   go build -o gopix
   ```

After these steps, the `gopix` binary will be available in your current project directory, ready to be used for image conversion tasks or moved to a different location to use it globally .

---


## 💻 Technologies Used

| Technology            | Description                                                                 |
|------------------------|-----------------------------------------------------------------------------|
| 🧠 **Golang**            | [go.dev](https://go.dev) — The core language powering GoPix: fast and efficient |
| 🛠️ **Cobra (CLI)**       | [spf13/cobra](https://github.com/spf13/cobra) — CLI commands, flags, and UX |
| 🎨 **Fatih/color**       | [fatih/color](https://github.com/fatih/color) — Terminal text styling and coloring |
| 🔄 **WebP encoder**      | [chai2010/webp](https://github.com/chai2010/webp) — Image conversion to/from WebP |
| 📏 **Resize**            | [nfnt/resize](https://github.com/nfnt/resize) — Image resizing utilities |
| 📉 **Progress bar**      | [schollz/progressbar](https://github.com/schollz/progressbar) — Beautiful terminal progress bar |
| 📦 **YAML config**       | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — Config file parser |
| 📜 **Logrus**            | [sirupsen/logrus](https://github.com/sirupsen/logrus) — Advanced logging framework |


---

## 🤝 Contributing

Contributions are welcome! Here’s how to get started:

1. Fork the repository  
2. Create a new branch:  
   `git checkout -b feature/YourFeature`  
3. Commit your changes:  
   `git commit -m "Add amazing feature"`  
4. Push to your branch:  
   `git push origin feature/YourFeature`  
5. Open a pull request

> 💡 Please open an issue first for major feature ideas or changes.

---

## 📄 License

This project is licensed under the **GPL-3.0 License**.  
See the [LICENSE](LICENSE) file for full details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/MostafaSensei106">MostafaSensei106</a>
</p>