package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"r34-client/controller"
	"r34-client/ui"
)

func main() {
	ctrl := controller.New()

	a := app.New()
	w := a.NewWindow("R34 Client")
	w.Resize(fyne.NewSize(800, 600))

	mainWindowContent := ui.NewMainWindowLayout(
		ui.SearchBar(func(query string) {
			ctrl.Search(query)
		}),
		ui.ImageList(ctrl.ListPostData),
		ui.Pagination(ctrl.TotalPage, ctrl.CurrentPage, func(page int) {
			ctrl.ChangePage(page)
		}),
	)
	// set all event handler or some shit

	w.SetContent(mainWindowContent)
	w.ShowAndRun()
}
