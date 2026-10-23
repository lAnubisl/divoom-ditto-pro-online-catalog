// Package rotation coordinates automatic playback and its schedule.
package rotation

import (
	"context"
	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/settings"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Runner struct {
	settings      *settings.Store
	selection     Selection
	source        Source
	display       Display
	mu            sync.Mutex
	changed       chan struct{}
	next          time.Time
	version       uint64
	started       bool
	running       bool
	lastError     string
	halted        bool
	cancel        context.CancelFunc
	takeover      bool
	restoreCancel context.CancelFunc
}

func New(prefs *settings.Store, selection Selection, source Source, display Display) *Runner {
	return &Runner{settings: prefs, selection: selection, source: source, display: display, changed: make(chan struct{}, 1), next: time.Now()}
}

func (r *Runner) State() State {
	r.mu.Lock()
	defer r.mu.Unlock()
	return State{Running: r.running, Error: r.lastError}
}

// Changed is called after a persisted settings change. Manual takeover cancels
// both the automatic upload and any pending clock restoration.
func (r *Runner) Changed(immediate, takeover bool) {
	r.mu.Lock()
	if immediate {
		r.next = time.Now()
		r.version++
		r.lastError = ""
		r.halted = false
	}
	prefs, err := r.settings.Read()
	if err == nil && !prefs.Enabled && r.cancel != nil {
		r.takeover = r.takeover || takeover
		r.cancel()
		if takeover && r.restoreCancel != nil {
			r.restoreCancel()
		}
	}
	r.mu.Unlock()
	select {
	case r.changed <- struct{}{}:
	default:
	}
}

func wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (r *Runner) restoreClock() error {
	r.mu.Lock()
	if r.takeover {
		r.mu.Unlock()
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 210*time.Second)
	r.restoreCancel = cancel
	r.mu.Unlock()
	defer func() { cancel(); r.mu.Lock(); r.restoreCancel = nil; r.mu.Unlock() }()
	prefs, err := r.settings.Read()
	if err == nil {
		err = r.display.ShowClock(ctx, prefs.Clock())
	}
	r.mu.Lock()
	takeover := r.takeover
	r.mu.Unlock()
	if takeover {
		return nil
	}
	if err != nil {
		return &clockRestoreError{err}
	}
	return nil
}

// nextAnimation learns unsupported IDs and immediately tries the next candidate.
// Failed retrieval never counts as playback history, including skipped artwork.
func (r *Runner) nextAnimation(ctx context.Context) (gallery.Item, []byte, error) {
	for {
		if err := ctx.Err(); err != nil {
			return gallery.Item{}, nil, err
		}
		item, err := r.selection.Next(ctx)
		if err != nil {
			return gallery.Item{}, nil, fmt.Errorf("select animation: %w", err)
		}
		data, err := r.source.Animation(ctx, item)
		if err == nil {
			return item, data, nil
		}
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		rejectErr := r.selection.Reject(cleanup, item)
		cancel()
		if rejectErr != nil {
			return gallery.Item{}, nil, errors.Join(fmt.Errorf("load animation: %w", err), fmt.Errorf("revert playback history: %w", rejectErr))
		}
		if !errors.Is(err, animation.ErrUnsupportedFormat) {
			return gallery.Item{}, nil, fmt.Errorf("load animation: %w", err)
		}
		log.Printf("Skipping animation %d: unsupported artwork format.", item.ID)
	}
}

func (r *Runner) cycle(ctx context.Context, duration time.Duration) (err error) {
	if err := r.display.Lock(ctx); err != nil {
		return err
	}
	defer r.display.Unlock()
	item, data, err := r.nextAnimation(ctx)
	if err != nil {
		return err
	}
	sent := false
	defer func() {
		if !sent {
			cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if historyErr := r.selection.Reject(cleanup, item); historyErr != nil {
				log.Printf("revert playback history: %v", historyErr)
			}
		}
	}()
	// This check prevents cancelled source work from starting a device request.
	if err := ctx.Err(); err != nil {
		return err
	}
	defer func() {
		if clockErr := r.restoreClock(); clockErr != nil {
			err = clockErr
		}
	}()
	if err := r.display.SendAnimation(ctx, data); err != nil {
		return fmt.Errorf("animation upload: %w", err)
	}
	sent = true
	log.Printf("Animation %s send completed; returning to clock in %s.", item.Filename(), duration)
	return wait(ctx, duration)
}

func (r *Runner) Start(ctx context.Context) (<-chan struct{}, error) {
	if _, err := r.settings.Read(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	if r.started {
		r.mu.Unlock()
		return nil, errors.New("automatic scheduler already started")
	}
	r.started = true
	r.mu.Unlock()
	done := make(chan struct{})
	go func() { defer close(done); r.run(ctx) }()
	return done, nil
}

func (r *Runner) run(ctx context.Context) {
	for {
		prefs, err := r.settings.Read()
		if err != nil {
			log.Printf("read rotation settings: %v", err)
			return
		}
		r.mu.Lock()
		halted := r.halted
		r.mu.Unlock()
		if !prefs.Enabled || halted {
			select {
			case <-ctx.Done():
				return
			case <-r.changed:
				continue
			}
		}
		r.mu.Lock()
		next := r.next
		r.mu.Unlock()
		timer := time.NewTimer(time.Until(next))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-r.changed:
			timer.Stop()
			continue
		case <-timer.C:
		}
		r.mu.Lock()
		// Re-read under the notification lock so a manual pause cannot be missed
		// between choosing a cycle and publishing its cancellation function.
		prefs, err = r.settings.Read()
		if err != nil || !prefs.Enabled || ctx.Err() != nil || time.Now().Before(r.next) {
			r.mu.Unlock()
			continue
		}
		interval, duration := prefs.Timing()
		next = r.next
		version := r.version
		cycleCtx, cancel := context.WithCancel(ctx)
		r.cancel = cancel
		r.takeover = false
		r.running = true
		r.mu.Unlock()
		err = r.cycle(cycleCtx, duration)
		cancelled := cycleCtx.Err() != nil
		cancel()
		r.mu.Lock()
		r.cancel = nil
		r.running = false
		var restoreErr *clockRestoreError
		if errors.As(err, &restoreErr) {
			r.halted = true
			if saveErr := r.settings.Update(func(v *settings.Values) { v.Enabled = false }); saveErr != nil {
				log.Printf("save stopped autoplay: %v", saveErr)
			}
			r.lastError = "Could not restore the clock. Autoplay has stopped; check the device."
		} else if err != nil && !cancelled {
			r.lastError = "Automatic send failed. We will retry in the next cycle."
		} else if err == nil {
			r.lastError = ""
		}
		if r.version == version {
			r.next = next.Add(interval)
			for !r.next.After(time.Now()) {
				r.next = r.next.Add(interval)
			}
		}
		r.mu.Unlock()
		if err != nil && !cancelled {
			log.Printf("Automatic cycle: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
