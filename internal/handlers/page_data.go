package handlers

import (
	"divoom-gallery-client/internal/gallery"
)

type pageData struct {
	Categories                     []gallery.Category
	Entries                        []gallery.Item
	Category, Page, Next, Previous int
	More                           bool
	CSRF, Error                    string
	Automatic                      automaticState
	Albums                         []gallery.Album
	Album                          int64
	Saved                          bool
}

func (d pageData) AnimationsPath() string {
	if d.Album > 0 {
		return "/albums/animations"
	}
	return "/categories/animations"
}
func (d pageData) QueryKey() string {
	if d.Album > 0 {
		return "album"
	}
	return "category"
}
func (d pageData) QueryID() int64 {
	if d.Album > 0 {
		return d.Album
	}
	return int64(d.Category)
}
