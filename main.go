package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
		}, func(q string) []string {
			return ctrl.GetAutoComplete(q)
		}),
		ui.ImageList(ctrl.ListPostData, func(id string) {
			srcUri := ctrl.ResolvePostIDDownloadURI(id)
			// return early
			if srcUri == nil {
				return
			}
			dlgSaveFile := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
				if err != nil {
					ctrl.SetStatusText(fmt.Sprintf("cancel save: %s", err))
					return
				}
				if closer == nil {
					return
				}
				ctrl.DownloadUri(srcUri, closer)
			}, w)
			dlgSaveFile.SetFileName(srcUri.Name())
			dlgSaveFile.Show()
		}),
		ui.Pagination(ctrl.TotalPage, ctrl.CurrentPage, func(page int) {
			ctrl.ChangePage(page)
		}),
		widget.NewLabelWithData(ctrl.StatusText),
	)
	w.SetMainMenu(ui.NewMainMenu(ctrl))
	// set all event handler or some shit

	w.SetContent(mainWindowContent)
	w.ShowAndRun()

	ctrl.OnClose()
}
