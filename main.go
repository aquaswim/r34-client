package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"r34-client/controller"
	"r34-client/ui"
	downloadList "r34-client/ui/download_list"
)

func main() {
	ctrl := controller.New()

	a := app.New()
	w := a.NewWindow("R34 Client")
	w.Resize(fyne.NewSize(800, 600))

	mainWindowContent := ui.NewMainWindowLayout(
		ui.SearchBar(func(query string) {
			ctrl.Search(query)
		}, func(q string) []string {
			return ctrl.GetAutoComplete(q)
		}),
		ui.ImageList(ctrl.ListPostData, func(id string) {
			ctrl.AddToDownloadList(id)
		}),
		ui.Pagination(ctrl.TotalPage, ctrl.CurrentPage, func(page int) {
			ctrl.ChangePage(page)
		}),
		widget.NewLabelWithData(ctrl.StatusText),
		downloadList.New(ctrl.DownloadPostIds, func() {
			if ctrl.IsDownloadListNotEmpty() {
				dlg := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
					ctrl.DownloadAllInList(uri)
				}, w)
				dlg.Show()
			}
		}),
	)
	w.SetMainMenu(ui.NewMainMenu(ctrl))
	// set all event handler or some shit

	w.SetContent(mainWindowContent)
	w.ShowAndRun()

	ctrl.OnClose()
}
