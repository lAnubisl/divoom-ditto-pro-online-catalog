package rotation

import "context"

type Metrics interface {
	WriteHistorySize(context.Context, int) error
}
