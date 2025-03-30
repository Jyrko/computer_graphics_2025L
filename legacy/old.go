package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)


const (
	BRIGHTNESS_FACTOR = 30    
	CONTRAST_FACTOR   = 1.5   
	GAMMA_FACTOR      = 1.8   
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


type Point struct {
	X, Y float32
}

type FunctionalFilter struct {
	Name   string
	Points []Point
}

var (
	
	IDENTITY_FILTER = FunctionalFilter{
		Name: "Identity",
		Points: []Point{
			{0, 0},
			{255, 255},
		},
	}

	currentFilter FunctionalFilter
	filterPoints  []Point
	filterCanvas  *canvas.Rectangle
	isDragging    bool
	selectedPoint int
)

var originalImage image.Image
var currentImage *image.RGBA

var imageCanvas *canvas.Image

func main() {
	a := app.NewWithID("computer-grahics.imagefilter")
	w := a.NewWindow("Image Filtering App")

	imageCanvas = canvas.NewImageFromImage(nil)
	imageCanvas.FillMode = canvas.ImageFillOriginal
	imageCanvas.SetMinSize(fyne.NewSize(200, 500))

	
	scroll := container.NewScroll(imageCanvas)
	scroll.SetMinSize(fyne.NewSize(800, 600)) 

	
	loadBtn := widget.NewButton("Load Image", func() {
		loadImage(w)
	})

	invertBtn := widget.NewButton("Invert", func() {
		if currentImage != nil {
			currentImage = invertImage(currentImage)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	brightnessBtn := widget.NewButton("Brightness", func() {
		if currentImage != nil {
			currentImage = brightnessCorrection(currentImage, BRIGHTNESS_FACTOR)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	contrastBtn := widget.NewButton("Contrast", func() {
		if currentImage != nil {
			currentImage = contrastEnhancement(currentImage, CONTRAST_FACTOR)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	gammaBtn := widget.NewButton("Gamma", func() {
		if currentImage != nil {
			currentImage = gammaCorrection(currentImage, GAMMA_FACTOR)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	blurBtn := widget.NewButton("Blur", func() {
		if currentImage != nil {
			currentImage = applyConvolution(currentImage, BLUR_KERNEL)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	gaussianBtn := widget.NewButton("Gaussian", func() {
		if currentImage != nil {
			currentImage = applyConvolution(currentImage, GAUSSIAN_KERNEL)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	sharpenBtn := widget.NewButton("Sharpen", func() {
		if currentImage != nil {
			currentImage = applyConvolution(currentImage, SHARPEN_KERNEL)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	edgeBtn := widget.NewButton("Edge Detect", func() {
		if currentImage != nil {
			currentImage = applyConvolution(currentImage, EDGE_DETECT_KERNEL)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	embossBtn := widget.NewButton("Emboss", func() {
		if currentImage != nil {
			currentImage = applyConvolution(currentImage, EMBOSS_KERNEL)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})


	dilateBtn := widget.NewButton("Dilation", func() {
			if currentImage != nil {
					currentImage = dilateImage(currentImage)
					imageCanvas.Image = currentImage
					imageCanvas.Refresh()
			}
	})

	erodeBtn := widget.NewButton("Erosion", func() {
			if currentImage != nil {
					currentImage = erodeImage(currentImage)
					imageCanvas.Image = currentImage
					imageCanvas.Refresh()
			}
	})

	resetBtn := widget.NewButton("Reset", func() {
		if originalImage != nil {
			currentImage = toRGBA(originalImage)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	saveBtn := widget.NewButton("Save", func() {
		if currentImage != nil {
			saveImage(w)
		}
	})

	
	filterCanvas = canvas.NewRectangle(color.White)
	filterCanvas.Resize(fyne.NewSize(256, 256))

	
	filterOverlay := newFilterOverlay()
	filterOverlay.Resize(fyne.NewSize(256, 256))

	
	filterStack := container.NewStack(filterCanvas, filterOverlay)

	
	currentFilter = IDENTITY_FILTER
	filterPoints = make([]Point, len(currentFilter.Points))
	copy(filterPoints, currentFilter.Points)

	
	applyFilterBtn := widget.NewButton("Apply Filter", func() {
		if currentImage != nil {
			currentImage = applyFunctionalFilter(currentImage, filterPoints)
			imageCanvas.Image = currentImage
			imageCanvas.Refresh()
		}
	})

	resetFilterBtn := widget.NewButton("Reset Filter", func() {
		filterPoints = make([]Point, len(IDENTITY_FILTER.Points))
		copy(filterPoints, IDENTITY_FILTER.Points)
		filterOverlay.Refresh()
	})

	
	filterEditor := container.NewVBox(
		filterStack,
		container.NewHBox(applyFilterBtn, resetFilterBtn),
	)

	
	content := container.NewHSplit(
		container.NewVBox(
				container.NewHBox(loadBtn, saveBtn, resetBtn),
				container.NewHBox(invertBtn, brightnessBtn, contrastBtn, gammaBtn),
				container.NewHBox(blurBtn, gaussianBtn, sharpenBtn, edgeBtn, embossBtn),
				container.NewHBox(dilateBtn, erodeBtn),
				scroll,
		),
		filterEditor,
)

	w.SetContent(content)
	
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}


func loadImage(win fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			return
		}

		
		img, _, err := image.Decode(reader)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		
		originalImage = img
		currentImage = toRGBA(img)
		imageCanvas.Image = currentImage
		imageCanvas.Refresh()

		
		
		bounds := currentImage.Bounds()
		imgWidth := float32(bounds.Dx())
		imgHeight := float32(bounds.Dy())

		
		screenSize := win.Canvas().Size()
		maxWidth := screenSize.Width
		maxHeight := screenSize.Height

		
		if imgWidth > maxWidth {
			imgWidth = maxWidth
		}
		
		if imgHeight+100 > maxHeight {
			imgHeight = maxHeight - 100
		}

		
		win.Resize(fyne.NewSize(imgWidth, imgHeight+100))

	}, win)
}


func saveImage(win fyne.Window) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if writer == nil {
			return
		}
		
		err = png.Encode(writer, currentImage)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
	}, win)
}


func toRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, src, bounds.Min, draw.Src)
	return rgba
}



func invertImage(src *image.RGBA) *image.RGBA {
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


func brightnessCorrection(src *image.RGBA, factor int) *image.RGBA {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			r, g, b, a := originalColor.RGBA()
			
			
			r8 := int(r>>8) + factor
			g8 := int(g>>8) + factor
			b8 := int(b>>8) + factor
			
			
			r8 = clamp(r8, 0, 255)
			g8 = clamp(g8, 0, 255)
			b8 = clamp(b8, 0, 255)
			
			result.Set(x, y, color.RGBA{uint8(r8), uint8(g8), uint8(b8), uint8(a>>8)})
		}
	}
	return result
}


func contrastEnhancement(src *image.RGBA, factor float64) *image.RGBA {
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
				uint8(clamp(int(r8), 0, 255)),
				uint8(clamp(int(g8), 0, 255)),
				uint8(clamp(int(b8), 0, 255)),
				uint8(a>>8),
			})
		}
	}
	return result
}


func gammaCorrection(src *image.RGBA, gamma float64) *image.RGBA {
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
				uint8(clamp(int(r32*255), 0, 255)),
				uint8(clamp(int(g32*255), 0, 255)),
				uint8(clamp(int(b32*255), 0, 255)),
				uint8(a>>8),
			})
		}
	}
	return result
}


func applyConvolution(src *image.RGBA, kernel [][]float64) *image.RGBA {
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
				R: uint8(clamp(int(r), 0, 255)),
				G: uint8(clamp(int(g), 0, 255)),
				B: uint8(clamp(int(b), 0, 255)),
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


func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}


func drawFilterGraph(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	
	
	drawGrid(img, w, h)
	
	
	if len(filterPoints) > 1 {
		for i := 0; i < len(filterPoints)-1; i++ {
			drawLine(img, filterPoints[i], filterPoints[i+1])
		}
	}
	
	
	for _, p := range filterPoints {
		drawPoint(img, p)
	}
	
	return img
}


func handleFilterClick(x, y float32) {
	
	for i, p := range filterPoints {
		if dist(x, y, p.X, p.Y) < 5 {
			isDragging = true
			selectedPoint = i
			return
		}
	}

	
	newPoint := Point{x, y}
	if isValidNewPoint(newPoint) {
		insertPoint(newPoint)
	}
}


func movePoint(index int, x, y float32) {
	if index <= 0 || index >= len(filterPoints)-1 {
		
		filterPoints[index].Y = clampFloat(y, 0, 255)
	} else {
		
		prevX := filterPoints[index-1].X
		nextX := filterPoints[index+1].X
		newX := clampFloat(x, prevX, nextX)
		filterPoints[index].X = newX
		filterPoints[index].Y = clampFloat(y, 0, 255)
	}
}


func applyFunctionalFilter(src *image.RGBA, points []Point) *image.RGBA {
	
	lut := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		lut[i] = uint8(interpolateY(float32(i), points))
	}

	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			result.Set(x, y, color.RGBA{
				R: lut[c.R],
				G: lut[c.G],
				B: lut[c.B],
				A: c.A,
			})
		}
	}
	
	return result
}


func clampFloat(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func dist(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func interpolateY(x float32, points []Point) float32 {
    if len(points) < 2 {
        return 0
    }

    
    if x <= points[0].X {
        return points[0].Y
    }
    if x >= points[len(points)-1].X {
        return points[len(points)-1].Y
    }

    
    var i int
    for i = 0; i < len(points)-1; i++ {
        if points[i].X <= x && points[i+1].X >= x {
            break
        }
    }
    
    
    p1 := points[i]
    p2 := points[i+1]
    
    
    if p2.X == p1.X {
        return p1.Y
    }
    
    t := (x - p1.X) / (p2.X - p1.X)
    return p1.Y + t*(p2.Y-p1.Y)
}

func isValidNewPoint(p Point) bool {
	
	for i := 0; i < len(filterPoints)-1; i++ {
		if p.X > filterPoints[i].X && p.X < filterPoints[i+1].X {
			return true
		}
	}
	return false
}

func insertPoint(p Point) {
	
	pos := 0
	for i := 0; i < len(filterPoints); i++ {
		if filterPoints[i].X > p.X {
			pos = i
			break
		}
	}
	
	
	filterPoints = append(filterPoints[:pos], append([]Point{p}, filterPoints[pos:]...)...)
}


func drawGrid(img *image.RGBA, w, h int) {
    
    bounds := img.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            img.Set(x, y, color.RGBA{250, 250, 250, 255})
        }
    }

    
    majorGridColor := color.RGBA{180, 180, 180, 255}
    minorGridColor := color.RGBA{220, 220, 220, 255}

    
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

    
    axisColor := color.RGBA{100, 100, 100, 255}
    
    for py := 0; py < h; py++ {
        img.Set(0, py, axisColor)
    }
    
    for px := 0; px < w; px++ {
        img.Set(px, h-1, axisColor)
    }
}

func drawLine(img *image.RGBA, p1, p2 Point) {
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

func drawPoint(img *image.RGBA, p Point) {
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



type FilterOverlay struct {
	canvas.Raster
}

func newFilterOverlay() *FilterOverlay {
	f := &FilterOverlay{}
	f.Raster = *canvas.NewRaster(drawFilterGraph)
	return f
}

func (f *FilterOverlay) MouseDown(ev *desktop.MouseEvent) {
	x, y := ev.Position.X, ev.Position.Y
	if ev.Button == desktop.MouseButtonPrimary {
		handleFilterClick(x, y)
	}
}

func (f *FilterOverlay) MouseUp(ev *desktop.MouseEvent) {
	isDragging = false
	selectedPoint = -1
}

func (f *FilterOverlay) MouseMoved(ev *desktop.MouseEvent) {
	if isDragging && selectedPoint >= 0 {
		movePoint(selectedPoint, ev.Position.X, ev.Position.Y)
		f.Refresh()
	}
}


func dilateImage(src *image.RGBA) *image.RGBA {
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

func erodeImage(src *image.RGBA) *image.RGBA {
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