package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"r34-client/entities"
	lazyImage "r34-client/ui/lazy_image"
)

func ImageList(data binding.UntypedList, onImgTap func(id string)) fyne.CanvasObject {
	return widget.NewGridWrapWithData(data, func() fyne.CanvasObject {
		return lazyImage.New()
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		i, err := item.(binding.Untyped).Get()
		if err != nil {
			log.Printf("error getting from untyped binding value: %s", err)
			return
		}
		post := i.(*entities.Post)
		li := object.(*lazyImage.LazyImage)

		li.ChangeImage(post)
		li.OnTap = func() {
			onImgTap(post.ID)
		}
	})
}
