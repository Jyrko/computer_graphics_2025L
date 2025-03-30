# Image Filter Editor

A desktop application for applying various image filters and effects, built with Go and Fyne GUI toolkit.

## Requirements

- Go 1.19 or later
- Fyne v2.x dependencies:
  - For macOS:
    - Xcode Command Line Tools: `xcode-select --install`
  - For Linux:
    - Required packages: `gcc libgl1-mesa-dev xorg-dev`
  - For Windows:
    - MinGW-w64 or MSYS2 with gcc
    - A C compiler (gcc)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/image-filter-editor.git
cd image-filter-editor
```

2. Install dependencies:
```bash
go mod download
```

## Running the Application

From the project root directory:

```bash
# Run directly
go run cmd/imagefilter/main.go

# Or build and run
go build -o imagefilter cmd/imagefilter/main.go
./imagefilter  # or imagefilter.exe on Windows
```



## Project Structure

```
.
├── cmd/
│   └── imagefilter/       # Application entry point
├── internal/
│   ├── filters/          # Image processing algorithms
│   │   ├── basic.go     # Basic filters (brightness, contrast, etc.)
│   │   └── quantize.go  # Dithering and quantization
│   ├── gui/             # User interface components
│   │   ├── window.go    # Main window implementation
│   │   └── overlay.go   # Filter controls overlay
│   └── utils/           # Helper functions
│       └── image.go     # Image conversion utilities
└── README.md
```

## Technologies Used

- [Go](https://golang.org/) - Programming language (1.19+)
- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit (v2.x)
- Standard library packages:
  - `image` - Core image processing
  - `image/color` - Color manipulation
  - `image/draw` - Drawing operations
  - `image/png` - Image format support

