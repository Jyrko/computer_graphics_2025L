package utils

import (
	"image"
	"image/draw"
)

func ToRGBA(src image.Image) *image.RGBA {
    bounds := src.Bounds()
    rgba := image.NewRGBA(bounds)
    draw.Draw(rgba, bounds, src, bounds.Min, draw.Src)
    return rgba
}

func Clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

func ClampFloat(value, min, max float32) float32 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}