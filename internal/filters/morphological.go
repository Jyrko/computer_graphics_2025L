package filters

import (
	"image"
	"image/color"
)

func DilateImage(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)


	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
			for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
					var maxR, maxG, maxB uint8
					a := src.RGBAAt(x, y).A
				
					for j := -1; j <= 1; j++ {
							for i := -1; i <= 1; i++ {
									p := src.RGBAAt(x+i, y+j)
									if p.R > maxR {
											maxR = p.R
									}
									if p.G > maxG {
											maxG = p.G
									}
									if p.B > maxB {
											maxB = p.B
									}
							}
					}
					result.Set(x, y, color.RGBA{maxR, maxG, maxB, a})
			}
	}
	copyBorder(src, result)
	return result
}

func ErodeImage(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)


	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
			for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
					minR, minG, minB := uint8(255), uint8(255), uint8(255)
					a := src.RGBAAt(x, y).A
				
					for j := -1; j <= 1; j++ {
							for i := -1; i <= 1; i++ {
									p := src.RGBAAt(x+i, y+j)
									if p.R < minR {
											minR = p.R
									}
									if p.G < minG {
											minG = p.G
									}
									if p.B < minB {
											minB = p.B
									}
							}
					}
					result.Set(x, y, color.RGBA{minR, minG, minB, a})
			}
	}
	copyBorder(src, result)
	return result
}

func copyBorder(src, dst *image.RGBA) {
	bounds := src.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, bounds.Min.Y, src.At(x, bounds.Min.Y))
			dst.Set(x, bounds.Max.Y-1, src.At(x, bounds.Max.Y-1))
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			dst.Set(bounds.Min.X, y, src.At(bounds.Min.X, y))
			dst.Set(bounds.Max.X-1, y, src.At(bounds.Max.X-1, y))
	}
}
