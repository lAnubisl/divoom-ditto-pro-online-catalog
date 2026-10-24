// Package playback coordinates manual commands and changes to automatic playback.
package playback

import (
	"context"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/settings"
	"time"
)

// Controller shares one command gate across manual playback and preference changes.
// Automatic cycles share the display lock with these commands.
type Controller struct {
	settings *settings.Store
	runner   Runner
	source   Source
	display  Display
	gate     chan struct{}
}

func New(prefs *settings.Store, runner Runner, source Source, display Display) *Controller {
	return &Controller{settings: prefs, runner: runner, source: source, display: display, gate: make(chan struct{}, 1)}
}

func (c *Controller) command(action func() error) error {
	select {
	case c.gate <- struct{}{}:
		defer func() { <-c.gate }()
		return action()
	default:
		return ErrBusy
	}
}

func (c *Controller) update(change func(*settings.Values), immediate, takeover bool) error {
	if err := c.settings.Update(change); err != nil {
		return &Error{Stage: "settings", Err: err}
	}
	c.runner.Changed(immediate, takeover)
	return nil
}

func (c *Controller) pause(takeover bool) error {
	return c.update(func(v *settings.Values) { v.Enabled = false }, false, takeover)
}

func (c *Controller) Pause() error {
	return c.command(func() error { return c.pause(false) })
}

func (c *Controller) Resume() error {
	return c.command(func() error {
		return c.update(func(v *settings.Values) { v.Enabled = true }, true, false)
	})
}

func (c *Controller) SetSource(kind string, id int64, name string) error {
	return c.command(func() error {
		return c.update(func(v *settings.Values) {
			if v.Source != kind || v.SourceID != id {
				v.SourcePage = 0
				v.PreviousItemID = 0
			}
			v.Source, v.SourceID, v.SourceName = kind, id, name
		}, false, false)
	})
}

func (c *Controller) SetTiming(interval, duration time.Duration) error {
	return c.command(func() error {
		return c.update(func(v *settings.Values) { v.Interval = interval.String(); v.Duration = duration.String() }, true, false)
	})
}

func (c *Controller) manual(ctx context.Context, action func() error) error {
	return c.command(func() error {
		// Persist the pause and cancel automatic work before waiting for the device.
		// A failed command deliberately leaves playback paused.
		if err := c.pause(true); err != nil {
			return err
		}
		if err := c.display.Lock(ctx); err != nil {
			return &Error{Stage: "wait", Err: err}
		}
		defer c.display.Unlock()
		return action()
	})
}

func (c *Controller) Send(ctx context.Context, item gallery.Item) error {
	return c.manual(ctx, func() error {
		data, err := c.source.Animation(ctx, item)
		if err != nil {
			return &Error{Stage: "load", Err: err}
		}
		if err := c.display.SendAnimation(ctx, data); err != nil {
			return &Error{Stage: "send", Err: err}
		}
		return nil
	})
}

func (c *Controller) RestoreClock(ctx context.Context) error {
	return c.manual(ctx, func() error {
		prefs, err := c.settings.Read()
		if err == nil {
			err = c.display.ShowClock(ctx, prefs.Clock())
		}
		if err != nil {
			return &Error{Stage: "clock", Err: err}
		}
		return nil
	})
}

// ApplyClock reports whether the device accepted the appearance, even if saving fails.
func (c *Controller) ApplyClock(ctx context.Context, clock settings.Clock) (applied bool, err error) {
	if err := clock.Validate(); err != nil {
		return false, err
	}
	err = c.manual(ctx, func() error {
		if err := c.display.ShowClock(ctx, clock); err != nil {
			return &Error{Stage: "clock", Err: err}
		}
		applied = true
		if err := c.settings.Update(func(v *settings.Values) { v.StyleID = clock.Style; v.StyleColor = clock.Color }); err != nil {
			return &Error{Stage: "save clock", Err: err}
		}
		return nil
	})
	return applied, err
}
