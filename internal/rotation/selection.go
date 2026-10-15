package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
)

type Selection interface {
	Next(context.Context) (gallery.Item, error)
	Reject(context.Context, gallery.Item) error
}
