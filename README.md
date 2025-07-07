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
- Multi-format support: PNG, JPG, WebP, JPEG
- Parallel processing: Uses all CPU cores for maximum speed
- Real-time progress bar with ETA
- Smart resume for interrupted conversions
- Custom quality and compression settings

### 🛠️ Advanced Capabilities
- Batch processing for folders and subfolders
- Size and resolution limits
- Configuration profiles with YAML support
- Dry-run mode to preview changes
- Backup of originals
- Rate limiting to prevent system overload
- Detailed post-process stats and reporting

### 🛡️ Security & Reliability
- Path validation to prevent directory traversal
- Safe defaults and permission checking
- Disk space validation before starting jobs
- Robust error handling and auto-retry mechanism

---

## 🛠️ Installation

## 📦 Easy Install (Linux / macOS / Windows)

> [!IMPORTANT]
> macOS users need to build GoPix from source.
> Currently, there is no pre-built binary availablea for arm64 architecture.

Download the latest pre-built binary for your platform from the [Releases](https://github.com/MostafaSensei106/GoPix/releases) page.

### 🐧 Linux

```bash
# Example for Linux (x86_64)
mv gopix ~/.local/bin/
```

>ℹ️ Make sure `~/.local/bin` is in your `$PATH`.

---

### 🍎 macOS

```bash
# Example for macOS (untested)
mv gopix /usr/local/bin/
```
---

### 🪟 Windows

1. Download `gopix.exe` from the [Releases](https://github.com/MostafaSensei106/GoPix/releases) page.
2. Move it to a folder like `C:\GoPix\bin`.
3. Add that folder to your **System PATH**:

   * Open *System Properties* → *Environment Variables* → *Path* → *Edit* → *Add new*.

Then you can run it from any Command Prompt:

```powershell
gopix help
```
---
## 🏗️ Build from Source (Linux, Windows, macOS)

> ![📝 Note] 
> GoPix uses a `Makefile` to build and install the CLI tool.  
> Make sure you have the `make` utility `Go` and `git`  installed on your system.  
> The script may adjust environment-specific paths depending on your OS.

---

### 🔧 Step 1: Install `make` (if not already installed)

#### For **Arch Linux** and based distros:
```bash
sudo pacman -S base-devel
```

#### For **Debian / Ubuntu** and based distros:
```bash
sudo apt install build-essential
```

#### For **Fedora** and based distros:
```bash
sudo dnf install make
```

#### For **openSUSE** and based distros:
```bash
sudo zypper install make
```

#### For **macOS**:
```bash
xcode-select --install
```

#### For **Windows**:
- Option 1: Install [MSYS2](https://www.msys2.org/) [recommended]
- Option 2: Use [Git Bash](https://gitforwindows.org/) and run the following command:
  ```bash
  pacman -Syu
  pacman -S make
  ```

---

### ⚙️ Step 2: Clone and Build

```bash
git clone --depth 1 https://github.com/MostafaSensei106/GoPix.git
cd GoPix
make
```

---

### ✅ Result

- This will compile GoPix from source optmized for your os and cpu architecture and install it locally.
- The binary will be placed in your system's executable path (like `/usr/local/bin` on Linux/macOS).
- You can now run:

```bash
gopix help
```


---

### Upgrading

> ![📝 Note]
> To upgrade GoPix, make sure you have the required development tools installed:
> `go`, `make`, and `git`.

To upgrade GoPix to the latest version, simply run
```bash
gopix upgrade
```
## This will:
  - Clone or update the latest source from GitHub.
  - Rebuild the binary using your current platform and       architecture.
  - Replace the old version automatically.

## OR

  get the latest pre-built binary for your platform from [Releases](https://github.com/MostafaSensei106/GoPix/releases) page and follow <a href="#installation">Installation Instructions</a>.

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
- Make

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
