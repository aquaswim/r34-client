package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type searchBarLayout struct {
	padding float32
}

func (s searchBarLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height-s.MinSize(objects).Height)
	for _, o := range objects {
		size := o.Size()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width+s.padding, 0))
	}
}

func (s searchBarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.Size()

		w += childSize.Width + s.padding
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func SearchBar(onSearch func(query string)) fyne.CanvasObject {
	searchTextInput := widget.NewEntry()
	searchButton := widget.NewButton("Search", func() {
		onSearch(searchTextInput.Text)
	})

	searchTextInput.Resize(fyne.NewSize(128, searchTextInput.MinSize().Height))
	searchButton.Resize(searchButton.MinSize())

	c := container.New(&searchBarLayout{padding: 5}, searchTextInput, searchButton)
	return c
}
