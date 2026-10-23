package rotation

import "context"

// History tracks successfully sent animations independently of persistence.
type History interface {
	Contains(string) bool
	Record(context.Context, string) error
}
