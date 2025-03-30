package main

import (
	"image-filter-editor/internal/gui"

	"fyne.io/fyne/v2/app"
)

func main() {
    a := app.NewWithID("computer-graphics.imagefilter")
    window := gui.NewMainWindow(a)
    window.Show()
    a.Run()
}