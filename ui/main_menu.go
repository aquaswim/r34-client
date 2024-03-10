package ui

import (
	"fyne.io/fyne/v2"
	"r34-client/controller"
)

func NewMainMenu(ctrl *controller.Controller) *fyne.MainMenu {
	return fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem("Clear Cache", func() {
		ctrl.ClearCache()
	})))
}
