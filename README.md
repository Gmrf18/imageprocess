# WebP Image Optimizer

A command-line tool written in Go that batch processes images: resizes them to multiple widths and converts them to **WebP** or **JPEG**, maintaining the original aspect ratio.

## What it does

- Reads all images from an input directory
- Detects orientation (portrait / landscape)
- Resizes to one or multiple predefined widths using Lanczos filter
- Exports to WebP or JPEG with chosen quality level
- Supports input formats: `.jpg`, `.jpeg`, `.png`, `.webp`, `.cr3`
- For `.cr3` files (Canon RAW) extracts embedded JPEG via `exiftool`
- Names output files as `name_orientation_width.ext`

## Usage flow

When running the program, answer the following interactive questions:

### 1. Directories

```
📁 Input directory [default: input]:
📁 Output directory [default: output]:
```

Press Enter to use the default values (`input/` and `output/`).

### 2. Output sizes

```
1) 320px  (Mobile S)
2) 468px  (Mobile L)
3) 768px  (Tablet)
4) 1080px (Full HD)
5) 1440px (2K)
6) 2160px (4K)
7) Custom (size in px)
```

You can choose one or multiple separated by comma: `1,3,5`  
If you choose `7`, enter the widths manually: `800, 1200`

> If the requested width exceeds the original image width, it is limited to the original to avoid pixelation.

### 3. Quality

```
1) Lossless / Minimal loss
2) Maximum Quality (95%)
3) High Quality    (80%) — Default
4) Medium-High Quality (65%)
5) Medium Quality  (45%)
6) Low Quality     (25%)
```

### 4. Output format

```
1) WebP  — Default
2) JPEG
```

> JPEG does not support lossless mode; if both options are combined, quality is adjusted to 100.

### 5. Processing

The program traverses the input directory and generates the requested files in the output directory for each image.

Example of generated files:

```
photo_landscape_768.webp
photo_landscape_1080.webp
photo_landscape_1440.webp
```

---

## Compilation

### Linux / WSL

```bash
go build -o imageprocess main.go
```

### Windows (.exe) from Linux / WSL

Requires the `mingw-w64` cross-compiler:

```bash
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o imageprocess.exe main.go
```

---

## Dependencies

| Package | Use |
|---|---|
| `github.com/chai2010/webp` | WebP encoding and decoding (CGO) |
| `github.com/disintegration/imaging` | Image resizing and opening |
| `exiftool` *(external)* | Embedded JPEG extraction from CR3 files |
