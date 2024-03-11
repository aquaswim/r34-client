package download_list

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type l struct {
	width float32
}

func (l l) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	btn := objects[1]
	list := objects[0]

	btn.Resize(fyne.NewSize(l.width, btn.MinSize().Height))
	list.Move(fyne.NewPos(0, 0))
	list.Resize(fyne.NewSize(l.width, containerSize.Height-btn.MinSize().Height))
	btn.Move(fyne.NewPos(0, list.Size().Height))

}

func (l l) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		w += childSize.Width
		h += childSize.Height
	}

	return fyne.NewSize(w, h)
}

func newLayout(width float32, list, btn fyne.CanvasObject) *fyne.Container {
	return container.New(&l{width: width}, list, btn)
}
