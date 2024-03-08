package lazy_image

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"r34-client/entities"
)

type LazyImage struct {
	widget.BaseWidget
	Image *canvas.Image
	Label *widget.Label
	Url   string
	Type  *entities.PostType
}

func New() *LazyImage {
	img := &LazyImage{
		Image: canvas.NewImageFromResource(theme.BrokenImageIcon()),
		Label: widget.NewLabel(""),
	}
	img.Image.SetMinSize(fyne.NewSize(120, 120))
	img.ExtendBaseWidget(img)
	return img
}

func NewFromPost(post *entities.Post) *LazyImage {
	img := New()
	go img.ChangeImage(post)
	return img
}

func (l *LazyImage) CreateRenderer() fyne.WidgetRenderer {
	//return widget.NewSimpleRenderer(l.Image)
	return widget.NewSimpleRenderer(container.NewStack(l.Image, l.Label))
}

func (l *LazyImage) Refresh() {
	l.Image.Refresh()
	l.Label.Refresh()
	l.BaseWidget.Refresh()
}

func (l *LazyImage) showErrorImage() {
	l.Image.Resource = theme.ErrorIcon()
	l.Label.Text = ""
}

func (l *LazyImage) ChangeImage(post *entities.Post) {
	url := post.ThumbnailURL
	// this function always being called don't know why
	if url == l.Url {
		return
	}
	defer l.Refresh()
	l.Url = url

	if url == "" {
		l.Image.File = ""
	}

	filePath, err := DownloadAndGetOutputFilePath(url)
	if err != nil {
		l.showErrorImage()
		log.Printf("parse image url: %s, err: %s", url, err)
		return
	}
	l.Image.Resource = nil
	l.Image.File = filePath

	switch post.Type {
	case entities.PostTypeImage:
		l.Label.Text = "IMG"
	case entities.PostTypeVideo:
		l.Label.Text = "VID"
	default:
		l.Label.Text = ""
	}
}
