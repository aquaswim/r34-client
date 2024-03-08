package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SearchBar(onSearch func(query string)) fyne.CanvasObject {
	searchTextInput := widget.NewEntry()
	searchButton := widget.NewButton("Search", func() {
		onSearch(searchTextInput.Text)
	})

	c := container.NewVBox(searchTextInput, searchButton)
	return c
}
