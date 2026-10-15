package source

import (
	"context"
	"divoom-gallery-client/internal/gallery"
)

type API interface {
	FetchCategoryPage(context.Context, int, int) (gallery.Page, error)
	FetchAlbumPage(context.Context, int64, int) (gallery.Page, error)
	FetchItem(context.Context, gallery.Item) ([]byte, error)
}
