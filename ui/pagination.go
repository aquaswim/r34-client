package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func Pagination(totalPage binding.Int, currentPage binding.Int, changePage func(page int)) fyne.CanvasObject {
	prevPageBtn := widget.NewButton("<<", func() {
		cp, _ := currentPage.Get()
		if cp < 1 {
			return
		}
		changePage(cp - 1)
	})
	nextPageBtn := widget.NewButton(">>", func() {
		cp, _ := currentPage.Get()
		tp, _ := totalPage.Get()
		if cp >= tp {
			changePage(tp)
		}
		changePage(cp + 1)
	})
	paginationLabel := widget.NewLabel("-/-")

	pagination := container.NewHBox(prevPageBtn, paginationLabel, nextPageBtn)

	updateListener := binding.NewDataListener(func() {
		_totalPage, _ := totalPage.Get()
		_currentPage, _ := currentPage.Get()
		paginationLabel.SetText(fmt.Sprintf("%d / %d", _currentPage, _totalPage))
		if _currentPage <= 1 {
			prevPageBtn.Hide()
		} else {
			prevPageBtn.Show()
		}
		if _currentPage >= _totalPage {
			nextPageBtn.Hide()
		} else {
			nextPageBtn.Show()
		}
	})

	totalPage.AddListener(updateListener)
	currentPage.AddListener(updateListener)

	return pagination
}
