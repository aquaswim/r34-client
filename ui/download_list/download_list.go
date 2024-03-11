package download_list

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DownloadList struct {
	widget.BaseWidget
	DownloadButton *widget.Button
	ClearButton    *widget.Button
	List           *widget.List
}

func New(postIds binding.StringList, onStartDownload func()) *DownloadList {
	d := &DownloadList{
		DownloadButton: widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
			onStartDownload()
		}),
		ClearButton: widget.NewButtonWithIcon("Clear", theme.CancelIcon(), func() {
			_ = postIds.Set([]string{})
		}),
		List: widget.NewListWithData(postIds, func() fyne.CanvasObject {
			return widget.NewLabel("-")
		}, func(i binding.DataItem, o fyne.CanvasObject) {
			str, _ := i.(binding.String).Get()
			o.(*widget.Label).SetText(str)
		}),
	}

	d.ExtendBaseWidget(d)
	return d
}

func (d *DownloadList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(newLayout(
		200,
		d.List,
		container.NewGridWithRows(1, d.DownloadButton, d.ClearButton),
	))
}
