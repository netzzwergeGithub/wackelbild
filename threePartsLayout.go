package main

import (
	"fyne.io/fyne/v2"
)

type threePartsLayout struct{}

func (d *threePartsLayout) MinSize(items []fyne.CanvasObject) fyne.Size {
	total := fyne.NewSize(0, 0)
	if len(items) != 3 {
		return total
	}
	for _, item := range items {
		total.Add(item.MinSize())
	}
	return total
}

func (d *threePartsLayout) Layout(items []fyne.CanvasObject, size fyne.Size) {
	if len(items) != 3 {
		return
	}
	items[0].Resize(fyne.NewSize(size.Width/4, size.Height))
	items[0].Move(fyne.NewPos(0, 0))
	items[1].Resize(fyne.NewSize(size.Width/2, size.Height))
	items[1].Move(fyne.NewPos(size.Width/4, 0))

	items[2].Resize(fyne.NewSize(size.Width/4, size.Height))
	items[2].Move(fyne.NewPos(3*size.Width/4, 0))

}
