package gui

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

type Point struct {
    X, Y float32
}

type FilterOverlay struct {
    canvas.Raster
    isDragging    bool
    selectedPoint int
    points        []Point
}

func NewFilterOverlay() *FilterOverlay {
    f := &FilterOverlay{
        points: []Point{
            {0, 0},
            {255, 255},
        },
        selectedPoint: -1,
    }
    f.Raster = *canvas.NewRaster(f.drawFilterGraph)
    return f
}

func (f *FilterOverlay) MouseDown(ev *desktop.MouseEvent) {
    x, y := ev.Position.X, ev.Position.Y
    if ev.Button == desktop.MouseButtonPrimary {
        f.handleClick(x, y)
    }
}

func (f *FilterOverlay) MouseUp(ev *desktop.MouseEvent) {
    f.isDragging = false
    f.selectedPoint = -1
}

func (f *FilterOverlay) MouseMoved(ev *desktop.MouseEvent) {
    if f.isDragging && f.selectedPoint >= 0 {
        f.movePoint(f.selectedPoint, ev.Position.X, ev.Position.Y)
        f.Refresh()
    }
}

func (f *FilterOverlay) drawFilterGraph(w, h int) image.Image {
    img := image.NewRGBA(image.Rect(0, 0, w, h))
    
    // Draw background and grid
    f.drawGrid(img, w, h)
    
    // Draw lines between points
    if len(f.points) > 1 {
        for i := 0; i < len(f.points)-1; i++ {
            f.drawLine(img, f.points[i], f.points[i+1])
        }
    }
    
    // Draw points
    for _, p := range f.points {
        f.drawPoint(img, p)
    }
    
    return img
}

func (f *FilterOverlay) drawGrid(img *image.RGBA, w, h int) {
    // Fill background
    bounds := img.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            img.Set(x, y, color.RGBA{250, 250, 250, 255})
        }
    }

    // Draw grid lines
    majorGridColor := color.RGBA{180, 180, 180, 255}
    minorGridColor := color.RGBA{220, 220, 220, 255}

    // Minor grid lines
    for i := 0; i <= 255; i += 16 {
        x := i * w / 255
        y := h - (i * h / 255)
        for py := 0; py < h; py++ {
            img.Set(x, py, minorGridColor)
        }
        for px := 0; px < w; px++ {
            img.Set(px, y, minorGridColor)
        }
    }

    // Major grid lines
    for i := 0; i <= 255; i += 64 {
        x := i * w / 255
        y := h - (i * h / 255)
        for py := 0; py < h; py++ {
            img.Set(x, py, majorGridColor)
        }
        for px := 0; px < w; px++ {
            img.Set(px, y, majorGridColor)
        }
    }

    // Draw axes
    axisColor := color.RGBA{100, 100, 100, 255}
    for py := 0; py < h; py++ {
        img.Set(0, py, axisColor)
    }
    for px := 0; px < w; px++ {
        img.Set(px, h-1, axisColor)
    }
}

func (f *FilterOverlay) drawLine(img *image.RGBA, p1, p2 Point) {
    x1, y1 := int(p1.X), int(255-p1.Y)
    x2, y2 := int(p2.X), int(255-p2.Y)
    
    lineColor := color.RGBA{0, 120, 255, 255}
    
    dx := float64(x2 - x1)
    dy := float64(y2 - y1)
    length := math.Sqrt(dx*dx + dy*dy)
    
    if length < 1 {
        return
    }
    
    dx /= length
    dy /= length
    
    for t := 0.0; t <= length; t++ {
        x := int(float64(x1) + dx*t)
        y := int(float64(y1) + dy*t)
        if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
            img.Set(x, y, lineColor)
            img.Set(x+1, y, lineColor)
            img.Set(x, y+1, lineColor)
        }
    }
}

func (f *FilterOverlay) drawPoint(img *image.RGBA, p Point) {
    x, y := int(p.X), int(255-p.Y)
    pointColor := color.RGBA{255, 50, 50, 255}
    outlineColor := color.RGBA{200, 0, 0, 255}
    
    radius := 4
    for dy := -radius; dy <= radius; dy++ {
        for dx := -radius; dx <= radius; dx++ {
            if dx*dx+dy*dy <= radius*radius {
                px := x + dx
                py := y + dy
                if px >= 0 && px < img.Bounds().Max.X && py >= 0 && py < img.Bounds().Max.Y {
                    if dx*dx+dy*dy == radius*radius {
                        img.Set(px, py, outlineColor)
                    } else {
                        img.Set(px, py, pointColor)
                    }
                }
            }
        }
    }
}

func (f *FilterOverlay) handleClick(x, y float32) {
    for i, p := range f.points {
        if f.dist(x, y, p.X, p.Y) < 5 {
            f.isDragging = true
            f.selectedPoint = i
            return
        }
    }

    newPoint := Point{x, y}
    if f.isValidNewPoint(newPoint) {
        f.insertPoint(newPoint)
    }
}

func (f *FilterOverlay) movePoint(index int, x, y float32) {
    if index <= 0 || index >= len(f.points)-1 {
        f.points[index].Y = f.clampFloat(y, 0, 255)
    } else {
        prevX := f.points[index-1].X
        nextX := f.points[index+1].X
        newX := f.clampFloat(x, prevX, nextX)
        f.points[index].X = newX
        f.points[index].Y = f.clampFloat(y, 0, 255)
    }
}

func (f *FilterOverlay) dist(x1, y1, x2, y2 float32) float32 {
    dx := x2 - x1
    dy := y2 - y1
    return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (f *FilterOverlay) clampFloat(value, min, max float32) float32 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

func (f *FilterOverlay) isValidNewPoint(p Point) bool {
    for i := 0; i < len(f.points)-1; i++ {
        if p.X > f.points[i].X && p.X < f.points[i+1].X {
            return true
        }
    }
    return false
}

func (f *FilterOverlay) insertPoint(p Point) {
    for i := 0; i < len(f.points)-1; i++ {
        if p.X > f.points[i].X && p.X < f.points[i+1].X {
            newPoints := make([]Point, 0, len(f.points)+1)
            newPoints = append(newPoints, f.points[:i+1]...)
            newPoints = append(newPoints, p)
            newPoints = append(newPoints, f.points[i+1:]...)
            f.points = newPoints
            break
        }
    }
    f.Refresh()
}

func (f *FilterOverlay) GetPoints() []Point {
    return f.points
}