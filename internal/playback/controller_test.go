package playback

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"divoom-gallery-client/internal/settings"
)

type settingsObserver struct {
	store *settings.Store
	seen  settings.Values
	calls int
}

func (r *settingsObserver) Changed(_, _ bool) {
	r.seen, _ = r.store.Read()
	r.calls++
}

type clockDisplay struct{}

func (clockDisplay) Lock(context.Context) error                      { return nil }
func (clockDisplay) Unlock()                                         {}
func (clockDisplay) SendAnimation(context.Context, []byte) error     { return nil }
func (clockDisplay) ShowClock(context.Context, settings.Clock) error { return nil }

func controllerPreferences(t *testing.T) (*settings.Store, string, settings.Values) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "settings.json")
	initial := settings.Defaults(3, time.Minute, time.Second)
	initial.SourcePage, initial.PreviousItemID, initial.SourceName = 7, 71, "old name"
	store, err := settings.New(path, initial)
	if err != nil {
		t.Fatal(err)
	}
	return store, path, initial
}

func TestSourceChangeResetsPageWithSource(t *testing.T) {
	for _, tc := range []struct {
		name, kind string
		id         int64
		page       int
	}{
		{name: "category ID", kind: "category", id: 4, page: 0},
		{name: "source kind", kind: "album", id: 3, page: 0},
		{name: "same source renamed", kind: "category", id: 3, page: 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prefs, path, want := controllerPreferences(t)
			observer := &settingsObserver{store: prefs}
			controller := New(prefs, observer, nil, nil)
			if err := controller.SetSource(tc.kind, tc.id, "new name"); err != nil {
				t.Fatal(err)
			}
			want.Source, want.SourceID, want.SourceName, want.SourcePage = tc.kind, tc.id, "new name", tc.page
			if tc.page == 0 {
				want.PreviousItemID = 0
			}
			if observer.calls != 1 || observer.seen != want {
				t.Fatalf("notification saw partial change: %+v", observer)
			}
			reopened, err := settings.New(path, settings.Defaults(0, time.Minute, time.Second))
			if err != nil {
				t.Fatal(err)
			}
			if got, err := reopened.Read(); err != nil || got != want {
				t.Fatalf("disk=%+v err=%v, want %+v", got, err, want)
			}
		})
	}
}

func TestTimingAndClockPreserveSourcePage(t *testing.T) {
	prefs, _, want := controllerPreferences(t)
	observer := &settingsObserver{store: prefs}
	controller := New(prefs, observer, nil, clockDisplay{})
	if err := controller.SetTiming(2*time.Minute, 2*time.Second); err != nil {
		t.Fatal(err)
	}
	if observer.seen.SourcePage != 7 {
		t.Fatal("timing change reset progress")
	}
	applied, err := controller.ApplyClock(context.Background(), settings.Clock{Style: 2, Color: "#abcdef"})
	if err != nil || !applied {
		t.Fatalf("apply clock: applied=%v err=%v", applied, err)
	}
	got, err := prefs.Read()
	if err != nil {
		t.Fatal(err)
	}
	want.Interval, want.Duration, want.StyleID, want.StyleColor, want.Enabled = "2m0s", "2s", 2, "#ABCDEF", false
	if got != want {
		t.Fatalf("clock/timing changed unrelated preferences: %+v, want %+v", got, want)
	}
}

func TestSourceChangeSaveFailurePreservesSourceAndPage(t *testing.T) {
	prefs, path, want := controllerPreferences(t)
	observer := &settingsObserver{store: prefs}
	controller := New(prefs, observer, nil, nil)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := controller.SetSource("album", 9, "new album"); err == nil {
		t.Fatal("expected save failure")
	}
	if got, err := prefs.Read(); err != nil || got != want {
		t.Fatalf("failed change published: %+v %v", got, err)
	}
	if observer.calls != 0 {
		t.Fatal("failed save notified runner")
	}
}
