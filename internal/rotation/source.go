package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
)

type Source interface {
	Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error)
	Animation(context.Context, gallery.Item) ([]byte, error)
}
