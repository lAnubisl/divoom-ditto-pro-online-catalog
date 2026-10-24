package rotation

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/settings"
	"divoom-gallery-client/internal/source"
)

func newPreferences(t *testing.T) *settings.Store {
	t.Helper()
	prefs, err := settings.New(filepath.Join(t.TempDir(), "settings.json"), settings.Defaults(0, time.Minute, time.Second))
	if err != nil {
		t.Fatal(err)
	}
	return prefs
}

type testSource struct {
	items       []gallery.Item
	requests    int
	loads       []int64
	nextError   error
	onCandidate func()
	onLoad      func(int64) error
}

func (s *testSource) NextAnimation(ctx context.Context) ([]byte, error) {
	s.requests++
	if s.onCandidate != nil {
		s.onCandidate()
	}
	if s.nextError != nil {
		return nil, s.nextError
	}
	if len(s.items) == 0 {
		return nil, source.ErrEmptyGallery
	}
	item := s.items[0]
	s.items = append(s.items[1:], item)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.loads = append(s.loads, item.ID)
	if s.onLoad != nil {
		if err := s.onLoad(item.ID); err != nil {
			return nil, err
		}
	}
	return []byte("GIF"), nil
}

type testDisplay struct {
	failure      error
	restoreError error
	locked       bool
	restored     bool
	sends        int
	onSend       func()
}

func (d *testDisplay) Lock(context.Context) error { d.locked = true; return nil }
func (d *testDisplay) Unlock()                    { d.locked = false }
func (d *testDisplay) SendAnimation(context.Context, []byte) error {
	d.sends++
	if d.onSend != nil {
		d.onSend()
	}
	return d.failure
}
func (d *testDisplay) ShowClock(context.Context, settings.Clock) error {
	d.restored = true
	return d.restoreError
}

func TestRunnerUsesEachCandidateInAdapterOrder(t *testing.T) {
	s := &testSource{items: []gallery.Item{{ID: 1}, {ID: 2}}}
	d := &testDisplay{}
	r := New(newPreferences(t), s, d)
	for i := 0; i < 5; i++ {
		if err := r.cycle(context.Background(), 0); err != nil {
			t.Fatal(err)
		}
	}
	if s.requests != 5 || !reflect.DeepEqual(s.loads, []int64{1, 2, 1, 2, 1}) {
		t.Fatalf("adapter order changed: requests=%d loads=%v", s.requests, s.loads)
	}
	if d.sends != 5 || d.locked || !d.restored {
		t.Fatalf("display=%+v", d)
	}
}

func TestRunnerAdvancesOnNextCycleAfterFailure(t *testing.T) {
	for _, stage := range []string{"download", "decode", "upload"} {
		t.Run(stage, func(t *testing.T) {
			failure := errors.New("network timeout")
			if stage == "decode" {
				failure = animation.ErrUnsupportedFormat
			}
			s := &testSource{items: []gallery.Item{{ID: 1}, {ID: 2}}}
			d := &testDisplay{}
			if stage == "upload" {
				d.failure = failure
			} else {
				s.onLoad = func(id int64) error {
					if id == 1 {
						return failure
					}
					return nil
				}
			}
			r := New(newPreferences(t), s, d)
			if err := r.cycle(context.Background(), 0); !errors.Is(err, failure) {
				t.Fatalf("lost failure: %v", err)
			}
			if s.requests != 1 || !reflect.DeepEqual(s.loads, []int64{1}) {
				t.Fatalf("retried in failed iteration: requests=%d loads=%v", s.requests, s.loads)
			}
			if d.locked {
				t.Fatal("failed cycle kept display lock")
			}
			if stage == "upload" && !d.restored {
				t.Fatal("failed upload did not restore clock")
			}
			d.failure = nil
			if err := r.cycle(context.Background(), 0); err != nil {
				t.Fatal(err)
			}
			if s.requests != 2 || !reflect.DeepEqual(s.loads, []int64{1, 2}) {
				t.Fatalf("did not advance: requests=%d loads=%v", s.requests, s.loads)
			}
		})
	}
}

func TestRunnerReturnsSourceErrorsWithoutLoading(t *testing.T) {
	for _, failure := range []error{source.ErrEmptyGallery, errors.New("page failed")} {
		s := &testSource{nextError: failure}
		r := New(newPreferences(t), s, &testDisplay{})
		if err := r.cycle(context.Background(), 0); !errors.Is(err, failure) {
			t.Fatalf("lost source error: %v", err)
		}
		if s.requests != 1 || len(s.loads) != 0 {
			t.Fatalf("unexpected source work: %+v", s)
		}
	}
}

func TestRunnerCancellationStopsWorkAndRestoresClock(t *testing.T) {
	for _, stage := range []string{"before candidate", "after candidate", "during load", "after send", "manual takeover"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := &testSource{items: []gallery.Item{{ID: 1}, {ID: 2}}}
			d := &testDisplay{}
			r := New(newPreferences(t), s, d)
			switch stage {
			case "before candidate":
				cancel()
			case "after candidate":
				s.onCandidate = cancel
			case "during load":
				s.onLoad = func(int64) error { cancel(); return context.Canceled }
			case "after send":
				d.onSend = cancel
			case "manual takeover":
				d.onSend = func() { r.mu.Lock(); r.takeover = true; r.mu.Unlock(); cancel() }
			}
			if err := r.cycle(ctx, time.Hour); !errors.Is(err, context.Canceled) {
				t.Fatalf("lost cancellation: %v", err)
			}
			if d.locked {
				t.Fatal("cancelled cycle retained lock")
			}
			sent := stage == "after send" || stage == "manual takeover"
			if (d.sends != 0) != sent {
				t.Fatalf("unexpected send count: %d", d.sends)
			}
			if d.restored != (stage == "after send") {
				t.Fatalf("unexpected restoration: %v", d.restored)
			}
		})
	}
}

func TestRunnerClockRestoreError(t *testing.T) {
	failure := errors.New("clock failed")
	s := &testSource{items: []gallery.Item{{ID: 1}}}
	d := &testDisplay{restoreError: failure}
	r := New(newPreferences(t), s, d)
	var restoreErr *clockRestoreError
	err := r.cycle(context.Background(), 0)
	if !errors.As(err, &restoreErr) || !errors.Is(err, failure) {
		t.Fatalf("lost clock restoration failure: %v", err)
	}
	if d.locked {
		t.Fatal("failed restoration retained lock")
	}
}

func TestRunnerWaitsUntilNextScheduledCycleAfterUploadFailure(t *testing.T) {
	prefs := newPreferences(t)
	interval := 200 * time.Millisecond
	if err := prefs.Update(func(v *settings.Values) { v.Interval = interval.String(); v.Duration = "1ms" }); err != nil {
		t.Fatal(err)
	}
	attempts := make(chan time.Time, 3)
	s := &testSource{items: []gallery.Item{{ID: 1}, {ID: 2}}}
	d := &testDisplay{failure: errors.New("upload failed")}
	d.onSend = func() { attempts <- time.Now() }
	r := New(prefs, s, d)
	ctx, cancel := context.WithCancel(context.Background())
	done, err := r.Start(ctx)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	defer func() { cancel(); <-done }()
	nextAttempt := func() time.Time {
		t.Helper()
		select {
		case at := <-attempts:
			return at
		case <-time.After(3 * time.Second):
			t.Fatal("scheduled upload did not occur")
			return time.Time{}
		}
	}
	first, second := nextAttempt(), nextAttempt()
	// Allow scheduling jitter while rejecting an immediate same-cycle retry.
	if second.Sub(first) < interval/2 {
		t.Fatalf("retried before scheduled delay: %s", second.Sub(first))
	}
	cancel()
	<-done
	if !reflect.DeepEqual(s.loads, []int64{1, 2}) {
		t.Fatalf("unexpected candidates: %v", s.loads)
	}
}
