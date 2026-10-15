package playback

import (
	"context"
	"divoom-gallery-client/internal/settings"
)

type Display interface {
	Lock(context.Context) error
	Unlock()
	SendAnimation(context.Context, []byte) error
	ShowClock(context.Context, settings.Clock) error
}
