package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"io"
	"r34-client/entities"
)

func (c *Controller) getPostIdsDetails(postIds []string) ([]*entities.PostDetail, error) {
	out := make([]*entities.PostDetail, 0, len(postIds))
	for _, id := range postIds {
		postDetail, err := c.dataSource.GetPostByID(id)
		if err != nil {
			return nil, err
		}
		out = append(out, postDetail)
		c.SetStatusText(fmt.Sprintf("postid: %s detail received", id))
	}

	return out, nil
}

func (c *Controller) downloadUri(srcUri fyne.URI, writer fyne.URIWriteCloser) {
	defer writer.Close()
	write, err := storage.CanWrite(writer.URI())
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error %s getting info writable for file: %s", writer.URI(), err))
		return
	}
	if write != true {
		c.SetStatusText(fmt.Sprintf("file %s not writeable", writer.URI()))
		return
	}

	c.SetStatusText(fmt.Sprintf("saving %s to %s", srcUri.Name(), writer.URI()))
	reader, err := storage.Reader(srcUri)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error reading %s: %s", reader.URI(), err))
		return
	}
	defer reader.Close()
	_, err = io.Copy(writer, reader)
	if err != nil {
		c.SetStatusText(fmt.Sprintf("error saving %s to %s", reader.URI(), writer.URI()))
		return
	}
	c.SetStatusText(fmt.Sprintf("%s saved to: %s", srcUri.Name(), writer.URI()))
}
