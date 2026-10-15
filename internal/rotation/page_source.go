package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
)

type PageSource interface {
	Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error)
}
