package main

import (
	"fyne.io/fyne/v2"
)

type twoPartsLayout struct {
}

func (d *twoPartsLayout) MinSize(items []fyne.CanvasObject) fyne.Size {
	total := fyne.NewSize(0, 0)
	for _, item := range items {
		total.Add(item.MinSize())
	}
	return total
}

func (d *twoPartsLayout) Layout(items []fyne.CanvasObject, size fyne.Size) {
	elementsCount := len(items)
	for i := 0; i < elementsCount-1; i++ {

		items[i].Resize(fyne.NewSize(size.Width/3, size.Height/float32(elementsCount-1)))
		items[i].Move(fyne.NewPos(0, float32(i)*size.Height/float32(elementsCount-1)))
	}
	items[len(items)-1].Resize(fyne.NewSize(2*size.Width/3, size.Height))
	items[len(items)-1].Move(fyne.NewPos(size.Width/3, 0))

}
