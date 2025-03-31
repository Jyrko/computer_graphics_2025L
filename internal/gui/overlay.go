package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FilterOverlay struct {
    container *fyne.Container
    sliders   map[string]*widget.Slider
    values    map[string]float64
    labels    map[string]*widget.Label
    onUpdate  func(string, float64)
}

func NewFilterOverlay() *FilterOverlay {
    f := &FilterOverlay{
        sliders: make(map[string]*widget.Slider),
        values:  make(map[string]float64),
        labels:  make(map[string]*widget.Label),
    }
    
    
    sliderConfigs := map[string]struct {
        min, max, value, step float64
        label string
    }{
        "brightness": {-100, 100, 0, 1, "Brightness"},
        "contrast":   {0, 3, 1, 0.1, "Contrast"},
        "gamma":      {0.1, 3, 1, 0.1, "Gamma"},
        "saturation": {0, 2, 1, 0.1, "Saturation"},
        "dither_levels": {2, 8, 2, 1, "Dither Levels"},
        "dither_size":   {2, 8, 2, 2, "Dither Map Size"},
        "num_colors":    {2, 256, 16, 1, "Number of Colors"},
    }

    var elements []fyne.CanvasObject
    
    
    for id, config := range sliderConfigs {
        nameLabel := widget.NewLabel(config.label)
        valueLabel := widget.NewLabel(formatValue(config.value))
        f.labels[id] = valueLabel
        
        slider := widget.NewSlider(config.min, config.max)
        slider.Step = config.step
        slider.Value = config.value
        f.values[id] = config.value
        f.sliders[id] = slider
        
        id := id
        slider.OnChanged = func(v float64) {
            f.values[id] = v
            f.labels[id].SetText(formatValue(v))
            if f.onUpdate != nil {
                f.onUpdate(id, v)
            }
        }
        

        sliderContainer := container.NewBorder(nil, nil, nil, valueLabel, slider)
        elements = append(elements, nameLabel, sliderContainer)
    }

    resetBtn := widget.NewButton("Reset All", func() {
        for id, config := range sliderConfigs {
            f.sliders[id].Value = config.value
            f.values[id] = config.value
            f.labels[id].SetText(formatValue(config.value))
            f.sliders[id].Refresh()
            if f.onUpdate != nil {
                f.onUpdate(id, config.value)
            }
        }
    })

    grayscaleBtn := widget.NewButton("Convert to Grayscale", func() {
        if f.onUpdate != nil {
            f.onUpdate("grayscale", 0)
        }
    })

    elements = append(elements, resetBtn)

    ditherBtn := widget.NewButton("Apply Dithering", func() {
        if f.onUpdate != nil {
            f.onUpdate("dither", f.values["dither_levels"])
        }
    })

    quantizeBtn := widget.NewButton("Quantize Colors", func() {
        if f.onUpdate != nil {
            f.onUpdate("quantize", f.values["num_colors"])
        }
    })

    elements = append(elements, grayscaleBtn, ditherBtn, quantizeBtn)

    f.container = container.NewVBox(elements...)
    
    return f
}

func (f *FilterOverlay) GetContainer() fyne.CanvasObject {
    return f.container
}

func (f *FilterOverlay) SetOnUpdate(callback func(param string, value float64)) {
    f.onUpdate = callback
}

func (f *FilterOverlay) GetValue(param string) float64 {
    return f.values[param]
}

func formatValue(value float64) string {
    if value == float64(int(value)) {
        return fmt.Sprintf("%.0f", value)
    }
    return fmt.Sprintf("%.1f", value)
}