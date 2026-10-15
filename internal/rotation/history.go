package rotation

import (
	"context"
)

type History interface {
	Load(context.Context) ([]string, error)
	Save(context.Context, []string) error
}
