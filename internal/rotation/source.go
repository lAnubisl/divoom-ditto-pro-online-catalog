package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
)

type Source interface {
	Animation(context.Context, gallery.Item) ([]byte, error)
}
