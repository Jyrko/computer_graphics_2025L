package utils

import (
	"image"
	"image/draw"
)

// ToRGBA converts any image.Image to *image.RGBA format
func ToRGBA(src image.Image) *image.RGBA {
    bounds := src.Bounds()
    rgba := image.NewRGBA(bounds)
    draw.Draw(rgba, bounds, src, bounds.Min, draw.Src)
    return rgba
}

// Clamp ensures an integer value stays within the given range
func Clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// ClampFloat ensures a float32 value stays within the given range
func ClampFloat(value, min, max float32) float32 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}