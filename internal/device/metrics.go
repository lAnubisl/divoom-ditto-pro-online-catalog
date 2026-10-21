package device

import "context"

type Metrics interface {
	WriteDeviceAPISuccess(context.Context, string, string) error
}
