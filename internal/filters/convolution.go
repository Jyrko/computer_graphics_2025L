package filters

import (
	"image"
	"image-filter-editor/internal/utils"
	"image/color"
)

var (
	BLUR_KERNEL = [][]float64{
		{1.0 / 9, 1.0 / 9, 1.0 / 9},
		{1.0 / 9, 1.0 / 9, 1.0 / 9},
		{1.0 / 9, 1.0 / 9, 1.0 / 9},
	}

	GAUSSIAN_KERNEL = [][]float64{
		{1.0 / 16, 2.0 / 16, 1.0 / 16},
		{2.0 / 16, 4.0 / 16, 2.0 / 16},
		{1.0 / 16, 2.0 / 16, 1.0 / 16},
	}

	SHARPEN_KERNEL = [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	EDGE_DETECT_KERNEL = [][]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}

	EMBOSS_KERNEL = [][]float64{
		{-2, -1, 0},
		{-1, 1, 1},
		{0, 1, 2},
	}
)


func ApplyConvolution(src *image.RGBA, kernel [][]float64) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	
	kernelSize := len(kernel)
	offset := kernelSize / 2

	for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
		for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
			var r, g, b float64
			
			
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					
					ix := x + (kx - offset)
					iy := y + (ky - offset)
					pixel := src.RGBAAt(ix, iy)
					
					
					k := kernel[ky][kx]
					r += float64(pixel.R) * k
					g += float64(pixel.G) * k
					b += float64(pixel.B) * k
				}
			}
			
			
			result.Set(x, y, color.RGBA{
				R: uint8(utils.Clamp(int(r), 0, 255)),
				G: uint8(utils.Clamp(int(g), 0, 255)),
				B: uint8(utils.Clamp(int(b), 0, 255)),
				A: src.RGBAAt(x, y).A,
			})
		}
	}

	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if y < bounds.Min.Y+offset || y >= bounds.Max.Y-offset ||
				x < bounds.Min.X+offset || x >= bounds.Max.X-offset {
				result.Set(x, y, src.At(x, y))
			}
		}
	}

	return result
}