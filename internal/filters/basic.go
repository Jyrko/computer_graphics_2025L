package filters

import (
	"image"
	"image-filter-editor/internal/utils"
	"image/color"
	"math"
	"sort"
)

const (
    BRIGHTNESS_FACTOR = 30    
    CONTRAST_FACTOR   = 1.5   
    GAMMA_FACTOR      = 1.8   
)

type Point struct {
	X, Y float64
}

func InvertImage(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			origColor := src.At(x, y)
			r, g, b, a := origColor.RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			
			invColor := color.RGBA{
				R: 255 - r8,
				G: 255 - g8,
				B: 255 - b8,
				A: a8,
			}
			result.Set(x, y, invColor)
		}
	}
	return result
}


func BrightnessCorrection(src *image.RGBA, factor int) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			r, g, b, a := originalColor.RGBA()
			
			
			r8 := int(r>>8) + factor
			g8 := int(g>>8) + factor
			b8 := int(b>>8) + factor
			
			
			r8 = utils.Clamp(r8, 0, 255)
			g8 = utils.Clamp(g8, 0, 255)
			b8 = utils.Clamp(b8, 0, 255)
			
			result.Set(x, y, color.RGBA{uint8(r8), uint8(g8), uint8(b8), uint8(a>>8)})
		}
	}
	return result
}


func ContrastEnhancement(src *image.RGBA, factor float64) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			r, g, b, a := originalColor.RGBA()
			
			r8 := float64(r>>8)
			g8 := float64(g>>8)
			b8 := float64(b>>8)
			
			
			r8 = ((r8 - 128) * factor) + 128
			g8 = ((g8 - 128) * factor) + 128
			b8 = ((b8 - 128) * factor) + 128
			
			
			result.Set(x, y, color.RGBA{
				uint8(utils.Clamp(int(r8), 0, 255)),
				uint8(utils.Clamp(int(g8), 0, 255)),
				uint8(utils.Clamp(int(b8), 0, 255)),
				uint8(a>>8),
			})
		}
	}
	return result
}


func GammaCorrection(src *image.RGBA, gamma float64) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			r, g, b, a := originalColor.RGBA()
			
			
			r32 := float64(r>>8) / 255.0
			g32 := float64(g>>8) / 255.0
			b32 := float64(b>>8) / 255.0
			
			
			r32 = math.Pow(r32, gamma)
			g32 = math.Pow(g32, gamma)
			b32 = math.Pow(b32, gamma)
			
			
			result.Set(x, y, color.RGBA{
				uint8(utils.Clamp(int(r32*255), 0, 255)),
				uint8(utils.Clamp(int(g32*255), 0, 255)),
				uint8(utils.Clamp(int(b32*255), 0, 255)),
				uint8(a>>8),
			})
		}
	}
	return result
}


// ApplyFunctionalFilter applies a custom transfer function defined by control points
func ApplyFunctionalFilter(src *image.RGBA, points []Point) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)

	// Sort points by X coordinate
	sort.Slice(points, func(i, j int) bool {
			return points[i].X < points[j].X
	})

	// Create lookup table
	var lookup [256]uint8
	for i := 0; i < 256; i++ {
			x := float64(i)
			// Find surrounding points
			var p1, p2 Point
			for j := 0; j < len(points)-1; j++ {
					if x >= points[j].X && x <= points[j+1].X {
							p1, p2 = points[j], points[j+1]
							break
					}
			}
			
			// Linear interpolation
			t := (x - p1.X) / (p2.X - p1.X)
			y := p1.Y + t*(p2.Y-p1.Y)
			
			// Clamp values
			lookup[i] = uint8(utils.Clamp(int(y), 0, 255))
	}

	// Apply lookup table to image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
					originalColor := src.At(x, y)
					r, g, b, a := originalColor.RGBA()
					
					result.Set(x, y, color.RGBA{
							R: lookup[r>>8],
							G: lookup[g>>8],
							B: lookup[b>>8],
							A: uint8(a>>8),
					})
			}
	}

	return result
}