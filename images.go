package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func getStartImages() (*canvas.Image, *canvas.Image) {
	initialImgLeft, _, err := image.Decode(bytes.NewReader(resourceJungfertigPng.Content()))
	if err != nil {
		fmt.Println(err)
	}
	initialImgRight, _, err := image.Decode(bytes.NewReader(resourceAltfertigPng.Content()))
	if err != nil {
		fmt.Println(err)
	}
	imageLeftWidget := canvas.NewImageFromImage(initialImgLeft)
	imageLeftWidget.FillMode = canvas.ImageFillContain

	imageRightWidget := canvas.NewImageFromImage(initialImgRight)
	imageRightWidget.FillMode = canvas.ImageFillContain
	return imageLeftWidget, imageRightWidget
}

func loadImage(imagePath string) (image.Image, error) {
	leftImageFile, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer leftImageFile.Close()
	imageLeft, _, err := image.Decode(leftImageFile)
	if err != nil {
		return nil, err
	}
	return imageLeft, nil
}

func getImageContainer(imgContainer *canvas.Image, addFileLoader bool, window fyne.Window, refresher func()) *fyne.Container {

	var top *fyne.Container

	if addFileLoader {
		openerFunction := func(read fyne.URIReadCloser, err error) {

			if read != nil {
				fmt.Println("User has chosen:", read.URI().String())
				img_loaded, err := loadImage(read.URI().Path())
				if err != nil {

					fmt.Println("Error loading image", read.URI().String(), err)
					return

				}
				imgContainer.Image = img_loaded
				imgContainer.Refresh()
				refresher()

			} else {
				fmt.Println("User caneled:", err)
			}

		}
		fileOpener := dialog.NewFileOpen(openerFunction, window)
		fileOpener.SetFilter(storage.NewMimeTypeFileFilter([]string{"image/jpeg", "image/png"}))

		openFile := widget.NewButtonWithIcon("Load Image", theme.FileIcon(), func() {
			fileOpener.Show()
		})
		top = container.NewCenter(openFile)

	}
	backgroundRect := canvas.NewRectangle(color.NRGBA{R: 125, G: 125, B: 125, A: 0xff})
	imgStacked := container.NewStack(backgroundRect, container.NewPadded(imgContainer))
	imageContainer := container.NewBorder(top, nil, nil, nil, imgStacked)
	return imageContainer

}
