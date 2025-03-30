package gui

import (
	"image"
	"image-filter-editor/internal/filters"
	"image-filter-editor/internal/utils"

	"image/color"
	"image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)
type MainWindow struct {
	window     fyne.Window
	image      *canvas.Image
	currentImg *image.RGBA
	origImg    image.Image
	filterCanvas  *canvas.Rectangle
	filterOverlay *FilterOverlay
	filterPoints  []filters.Point
}


func NewMainWindow(app fyne.App) *MainWindow {
    w := &MainWindow{
        window: app.NewWindow("Image Filtering App"),
				filterPoints: []filters.Point{ 
					{X: 0, Y: 0},
					{X: 255, Y: 255},
			},
    }

    w.image = canvas.NewImageFromImage(nil)
    w.image.FillMode = canvas.ImageFillOriginal
    w.image.SetMinSize(fyne.NewSize(200, 500))

    scroll := container.NewScroll(w.image)
    scroll.SetMinSize(fyne.NewSize(800, 600))

    buttons := w.createButtons()
    
    w.window.Resize(fyne.NewSize(800, 600))

		 w.filterCanvas = canvas.NewRectangle(color.White)
		 w.filterCanvas.Resize(fyne.NewSize(256, 256))
		 
		 w.filterOverlay = NewFilterOverlay()
    w.filterOverlay.SetOnUpdate(func(param string, value float64) {
        if w.currentImg == nil {
            return
        }
        
        switch param {
        case "brightness":
            w.currentImg = filters.BrightnessCorrection(w.currentImg, int(value))
        case "contrast":
            w.currentImg = filters.ContrastEnhancement(w.currentImg, value)
        case "gamma":
            w.currentImg = filters.GammaCorrection(w.currentImg, value)
        }
        
        w.image.Image = w.currentImg
        w.image.Refresh()
    })

    content := container.NewHSplit(
        container.NewVBox(buttons, scroll),
        w.filterOverlay.GetContainer(),
    )

    w.window.SetContent(content)
    return w
}

func (w *MainWindow) createButtons() *fyne.Container {
	loadBtn := widget.NewButton("Load Image", func() {
			loadImage(w)
	})

	saveBtn := widget.NewButton("Save", func() {
			if w.currentImg != nil {
					saveImage(w)
			}
	})

	resetBtn := widget.NewButton("Reset", func() {
			if w.origImg != nil {
					w.currentImg = utils.ToRGBA(w.origImg)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	invertBtn := widget.NewButton("Invert", func() {
			if w.currentImg != nil {
					w.currentImg = filters.InvertImage(w.currentImg)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	brightnessBtn := widget.NewButton("Brightness", func() {
			if w.currentImg != nil {
					w.currentImg = filters.BrightnessCorrection(w.currentImg, filters.BRIGHTNESS_FACTOR)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	contrastBtn := widget.NewButton("Contrast", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ContrastEnhancement(w.currentImg, filters.CONTRAST_FACTOR)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	gammaBtn := widget.NewButton("Gamma", func() {
			if w.currentImg != nil {
					w.currentImg = filters.GammaCorrection(w.currentImg, filters.GAMMA_FACTOR)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	blurBtn := widget.NewButton("Blur", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ApplyConvolution(w.currentImg, filters.BLUR_KERNEL)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	gaussianBtn := widget.NewButton("Gaussian", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ApplyConvolution(w.currentImg, filters.GAUSSIAN_KERNEL)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	sharpenBtn := widget.NewButton("Sharpen", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ApplyConvolution(w.currentImg, filters.SHARPEN_KERNEL)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	edgeBtn := widget.NewButton("Edge Detect", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ApplyConvolution(w.currentImg, filters.EDGE_DETECT_KERNEL)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	embossBtn := widget.NewButton("Emboss", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ApplyConvolution(w.currentImg, filters.EMBOSS_KERNEL)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	dilateBtn := widget.NewButton("Dilation", func() {
			if w.currentImg != nil {
					w.currentImg = filters.DilateImage(w.currentImg)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	erodeBtn := widget.NewButton("Erosion", func() {
			if w.currentImg != nil {
					w.currentImg = filters.ErodeImage(w.currentImg)
					w.image.Image = w.currentImg
					w.image.Refresh()
			}
	})

	return container.NewVBox(
			container.NewHBox(loadBtn, saveBtn, resetBtn),
			container.NewHBox(invertBtn, brightnessBtn, contrastBtn, gammaBtn),
			container.NewHBox(blurBtn, gaussianBtn, sharpenBtn, edgeBtn, embossBtn),
			container.NewHBox(dilateBtn, erodeBtn),
	)
}

func (w *MainWindow) Show() {
    w.window.ShowAndRun()
}

func loadImage(w *MainWindow) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
					dialog.ShowError(err, w.window)
					return
			}
			if reader == nil {
					return
			}

			img, _, err := image.Decode(reader)
			if err != nil {
					dialog.ShowError(err, w.window)
					return
			}

			w.origImg = img
			w.currentImg = utils.ToRGBA(img)
			w.image.Image = w.currentImg
			w.image.Refresh()

			bounds := w.currentImg.Bounds()
			imgWidth := float32(bounds.Dx())
			imgHeight := float32(bounds.Dy())

			screenSize := w.window.Canvas().Size()
			maxWidth := screenSize.Width
			maxHeight := screenSize.Height

			if imgWidth > maxWidth {
					imgWidth = maxWidth
			}
			if imgHeight+100 > maxHeight {
					imgHeight = maxHeight - 100
			}

			w.window.Resize(fyne.NewSize(imgWidth, imgHeight+100))
	}, w.window)
}

func saveImage(w *MainWindow) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
					dialog.ShowError(err, w.window)
					return
			}
			if writer == nil {
					return
			}
			
			err = png.Encode(writer, w.currentImg)
			if err != nil {
					dialog.ShowError(err, w.window)
					return
			}
			writer.Close()
	}, w.window)
}