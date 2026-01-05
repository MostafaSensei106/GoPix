# v2.0.0

This major release introduces a complete overhaul of the image processing engine and adds several powerful new features.

## ✨ New Features

- **New Image Formats**: Added support for modern and professional image formats:
  - `AVIF`
  - `HEIF` / `HEIC`
  - `TIFF`
  - `GIF`
- **Metadata Control**: You can now control how EXIF and other metadata is handled with the `--metadata` flag.
  - `keep` (default): Retains all metadata.
  - `strip`: Removes all metadata for smaller file sizes.
  - `strip-location` (Future): Will remove GPS data only.

## 🚀 Performance & Engine

- **New Engine**: Replaced the standard Go image libraries with `govips` (using `libvips`). This results in a **4-8x performance increase** and significantly lower memory usage, especially with large images.
- **Dependency**: `GoPix` now requires `libvips` to be installed on the system.

## 🛠️ Improvements & Fixes

- **Improved Rate Limiting**: The worker pool's rate-limiting logic is now more stable and efficient, using a blocking wait instead of a non-blocking check.
- **Better Error Reporting**: Failure reasons in the final report are now categorized (e.g., Corrupted Image, Permission Denied), making it easier to diagnose issues with large batches.
- **Updated Dependencies**: Added `github.com/davidbyttow/govips/v2`.
- **Removed Dependencies**: Removed `github.com/chai2010/webp` and `github.com/nfnt/resize` as their functionality is now covered by `govips`.
