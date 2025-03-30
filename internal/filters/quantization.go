package filters

import (
	"image"
	"image-filter-editor/internal/utils"
	"image/color"
	"sort"
)

func ToGrayscale(src *image.RGBA) *image.RGBA {
    bounds := src.Bounds()
    result := image.NewRGBA(bounds)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            originalColor := src.At(x, y)
            r, g, b, _ := originalColor.RGBA()
            
          
            gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256.0)
            result.Set(x, y, color.RGBA{gray, gray, gray, 255})
        }
    }
    return result
}

func OrderedDithering(src *image.RGBA, mapSize int, levels int) *image.RGBA {
    bounds := src.Bounds()
    result := image.NewRGBA(bounds)
    
  
    thresholdMap := makeThresholdMap(mapSize)
    
    step := 255.0 / float64(levels-1)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            originalColor := src.At(x, y)
            r, g, b, a := originalColor.RGBA()
            
          
            threshold := thresholdMap[(y%mapSize)][(x%mapSize)]
            
          
            newR := ditherValue(uint8(r>>8), threshold, step)
            newG := ditherValue(uint8(g>>8), threshold, step)
            newB := ditherValue(uint8(b>>8), threshold, step)
            
            result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a>>8)})
        }
    }
    return result
}

func PopularityQuantization(src *image.RGBA, numColors int) *image.RGBA {
    bounds := src.Bounds()
    result := image.NewRGBA(bounds)
    
  
    colorCount := make(map[color.RGBA]int)
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            c := src.RGBAAt(x, y)
            colorCount[c]++
        }
    }
    
  
    type colorFreq struct {
        color color.RGBA
        count int
    }
    
    var frequencies []colorFreq
    for c, count := range colorCount {
        frequencies = append(frequencies, colorFreq{c, count})
    }
    
    sort.Slice(frequencies, func(i, j int) bool {
        return frequencies[i].count > frequencies[j].count
    })
    
  
    palette := make([]color.RGBA, 0, numColors)
    for i := 0; i < numColors && i < len(frequencies); i++ {
        palette = append(palette, frequencies[i].color)
    }
    
  
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            original := src.RGBAAt(x, y)
            nearest := findNearestColor(original, palette)
            result.Set(x, y, nearest)
        }
    }
    
    return result
}

func makeThresholdMap(size int) [][]float64 {
  
    basic := [][]float64{
        {0.0, 0.5},
        {0.75, 0.25},
    }
    
    if size == 2 {
        return basic
    }
    
  
    result := make([][]float64, size)
    for i := range result {
        result[i] = make([]float64, size)
    }
    
  
    scale := float64(size * size)
    for y := 0; y < size; y++ {
        for x := 0; x < size; x++ {
            result[y][x] = float64(((x^y)*size + x) % (size*size)) / scale
        }
    }
    
    return result
}

func ditherValue(value uint8, threshold, step float64) uint8 {
    normalized := float64(value) / 255.0
    if normalized > threshold {
        level := int((normalized * 255.0) / step)
        return uint8(utils.Clamp(level*int(step), 0, 255))
    }
    level := int((normalized * 255.0) / step) - 1
    return uint8(utils.Clamp(level*int(step), 0, 255))
}

func findNearestColor(c color.RGBA, palette []color.RGBA) color.RGBA {
    minDist := float64(1<<32 - 1)
    var nearest color.RGBA
    
    for _, p := range palette {
        dist := colorDistance(c, p)
        if dist < minDist {
            minDist = dist
            nearest = p
        }
    }
    
    return nearest
}

func colorDistance(c1, c2 color.RGBA) float64 {
    dr := float64(c1.R) - float64(c2.R)
    dg := float64(c1.G) - float64(c2.G)
    db := float64(c1.B) - float64(c2.B)
    return dr*dr + dg*dg + db*db
}