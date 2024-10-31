package main

//go:generate fyne bundle -o bundled.go images/Altfertig.png
//go:generate fyne bundle -o bundled.go -append images/Jungfertig.png
//go:generate fyne bundle -o bundled.go -append images/open.png

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type configuration struct {
	viewAngle float64
	MAX_ANGLE float64
}

var myConfiguration = configuration{
	viewAngle: -60.,
	MAX_ANGLE: 120.,
}

func main() {

	myApp := app.NewWithID("de.netzzwerge.wackelbild.preferences")
	window := myApp.NewWindow("MaxLayout")

	// startImages laden
	imageLeftWidget, imageRightWidget := getStartImages()

	combineImgs := getCombined(imageLeftWidget.Image, imageRightWidget.Image, math.Abs(myConfiguration.viewAngle))

	backgroundRectCenter := canvas.NewRectangle(color.NRGBA{R: 205, G: 205, B: 205, A: 0xff})

	center_image := canvas.NewImageFromImage(combineImgs.compressed)
	center_image.FillMode = canvas.ImageFillContain

	centerImageContainer := container.NewPadded(container.NewStack(backgroundRectCenter, center_image))

	stripedButton := widget.NewButton("Striped", func() { center_image.Image = combineImgs.striped; center_image.Refresh() })
	compressedButton := widget.NewButton("Compressed", func() { center_image.Image = combineImgs.compressed; center_image.Refresh() })
	angle_binding := binding.BindFloat(&myConfiguration.viewAngle)
	angle_Text := widget.NewLabel("DUMMY")
	angle_binding.AddListener(binding.NewDataListener(func() {
		combineImgs = getCombined(imageLeftWidget.Image, imageRightWidget.Image, math.Abs(myConfiguration.viewAngle))
		center_image.Image = combineImgs.compressed
		center_image.Refresh()
		angle_Text.SetText(fmt.Sprintf("current angle: %03d", int(math.Abs(myConfiguration.viewAngle))))
	}))
	angle_binding.Reload()

	refeshCenterImage := func() {
		combineImgs = getCombined(imageLeftWidget.Image, imageRightWidget.Image, math.Abs(myConfiguration.viewAngle))
		center_image.Image = combineImgs.compressed
		center_image.Refresh()
	}

	leftImageContainer := getImageContainer(imageLeftWidget, true, window, refeshCenterImage)
	rightImageContainer := getImageContainer(imageRightWidget, true, window, refeshCenterImage)

	angleSlider := widget.NewSliderWithData(-120., 0., angle_binding)
	angleSlider.Orientation = 1
	paddedSlider := container.NewGridWithRows(2, angleSlider, container.NewVBox(angle_Text, stripedButton, compressedButton))

	// centerTopContainer := container.NewCenter(angle_Text)
	centerContainer := container.NewBorder(nil, nil, paddedSlider, nil, centerImageContainer)

	// backgroundRectRight := canvas.NewRectangle(color.NRGBA{R: 125, G: 125, B: 125, A: 0xff})
	// threeParts := container.New(&threePartsLayout{}, leftImageContainer, centerContainer, rightImageContainer)
	threeParts := container.New(&twoPartsLayout{}, leftImageContainer, rightImageContainer, centerContainer)
	backgroundRectGlobal := canvas.NewRectangle(color.NRGBA{R: 205, G: 205, B: 205, A: 0xff})
	myContainer := container.NewStack(backgroundRectGlobal, threeParts)
	window.SetContent(myContainer)
	window_x := float32(1000.)
	window_y := float32(600.)
	window.Resize(fyne.NewSize(window_x, window_y))
	window.ShowAndRun()

}
