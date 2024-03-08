package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewMainWindowLayout(
	searchBar fyne.CanvasObject,
	posts fyne.CanvasObject,
	pagination fyne.CanvasObject,
	statusText fyne.CanvasObject,
) fyne.CanvasObject {
	mainContainer := container.NewBorder(
		nil,
		nil,
		container.NewVBox(
			searchBar,
			// download list TBA
		),
		nil,
		container.NewBorder(
			nil,
			container.NewHBox(pagination, statusText),
			nil,
			nil,
			posts,
		),
	)
	return mainContainer
}
