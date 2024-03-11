package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
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

func SearchBar(onSearch func(query string), getAutoComplete func(q string) []string) fyne.CanvasObject {
	searchTextInput := NewCompletionEntry([]string{})
	searchTextInput.OnSubmitted = func(s string) {
		searchTextInput.Disable()
		searchTextInput.Options = []string{}
		defer searchTextInput.Enable()

		tags := strings.Split(searchTextInput.Text, " ")
		if len(tags) > 0 {
			searchQuery := tags[len(tags)-1]
			opt := getAutoComplete(searchQuery)
			if opt != nil {
				searchTextInput.Options = opt
			}
		}

		searchTextInput.ShowCompletion()
	}
	searchTextInput.SetPlaceHolder("enter to show autocomplete")
	searchTextInput.CustomSetText = func(selectedOpt string) string {
		tags := strings.Split(searchTextInput.Text, " ")
		if len(tags) > 0 {
			tags = tags[:len(tags)-1]
		}
		newTags := append(tags, selectedOpt)

		return strings.Join(newTags, " ") + " "
	}

	searchButton := widget.NewButton("Search", func() {
		onSearch(searchTextInput.Text)
	})

	searchTextInput.Resize(fyne.NewSize(300, searchTextInput.MinSize().Height))
	searchButton.Resize(searchButton.MinSize())

	c := container.New(&searchBarLayout{padding: 5}, searchTextInput, searchButton)
	return c
}
