package rotation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/history"
	"divoom-gallery-client/internal/settings"
)

type retryPages struct{ items []gallery.Item }

func (p retryPages) Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error) {
	return func(_ context.Context, page int) ([]gallery.Item, bool, error) {
		if page > 1 {
			return nil, true, nil
		}
		return p.items, true, nil
	}, nil
}

// Models the source's persistent filtering without requiring device access.
type retrySource struct {
	pages   *retryPages
	failure error
	calls   []int64
}

func (s *retrySource) Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error) {
	return s.pages.Pages()
}

func (s *retrySource) Animation(_ context.Context, item gallery.Item) ([]byte, error) {
	s.calls = append(s.calls, item.ID)
	if item.ID == 1 {
		if errors.Is(s.failure, animation.ErrUnsupportedFormat) {
			s.pages.items = s.pages.items[1:]
		}
		return nil, s.failure
	}
	return []byte("cached GIF"), nil
}

func TestRotationImmediatelyAdvancesPastUnsupportedItem(t *testing.T) {
	for _, allUnsupported := range []bool{false, true} {
		t.Run(map[bool]string{false: "next supported", true: "all unsupported"}[allUnsupported], func(t *testing.T) {
			store, err := newHistory(t, filepath.Join(t.TempDir(), "history.json"), nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			pages := &retryPages{items: []gallery.Item{{ID: 1}, {ID: 2}}}
			if allUnsupported {
				pages.items = pages.items[:1]
			}
			source := &retrySource{pages: pages, failure: animation.ErrUnsupportedFormat}
			runner := New(nil, store, source, nil)
			item, data, err := runner.nextAnimation(context.Background())
			if allUnsupported {
				if !errors.Is(err, ErrEmptyGallery) {
					t.Fatalf("expected empty gallery, got %v", err)
				}
			} else if err != nil || item.ID != 2 || len(data) == 0 || len(source.calls) != 2 {
				t.Fatalf("did not advance in same cycle: item=%v calls=%v err=%v", item, source.calls, err)
			}
			names, err := store.Load(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if len(names) != 0 {
				t.Fatalf("unsupported item entered history: %v", names)
			}
		})
	}
}

func TestRotationDoesNotSkipNetworkFailures(t *testing.T) {
	store, err := newHistory(t, filepath.Join(t.TempDir(), "history.json"), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	pages := &retryPages{items: []gallery.Item{{ID: 1}, {ID: 2}}}
	failure := errors.New("network timeout")
	source := &retrySource{pages: pages, failure: failure}
	runner := New(nil, store, source, nil)
	if _, _, err := runner.nextAnimation(context.Background()); !errors.Is(err, failure) {
		t.Fatalf("network error lost: %v", err)
	}
	if len(source.calls) != 1 || len(pages.items) != 2 {
		t.Fatal("network failure was treated as unsupported")
	}
	names, err := store.Load(context.Background())
	if err != nil || len(names) != 0 {
		t.Fatalf("failed download remained in history: %v, %v", names, err)
	}
}

type historySamples struct{ sizes []int }

func (m *historySamples) WriteHistorySize(_ context.Context, size int) error {
	m.sizes = append(m.sizes, size)
	return nil
}

func TestRunnerSelectionHistory(t *testing.T) {
	full := make([]string, history.Limit)
	for i := range full {
		full[i] = fmt.Sprintf("%d.gif", i+10)
	}
	for _, tc := range []struct {
		name        string
		initial     []string
		wantID      int64
		wantHistory []string
	}{
		{name: "skip previously played", initial: []string{"1.gif"}, wantID: 2, wantHistory: []string{"1.gif", "2.gif"}},
		{name: "restart exhausted source", initial: []string{"1.gif", "2.gif"}, wantID: 1, wantHistory: []string{"1.gif"}},
		{name: "limit retained history", initial: full, wantID: 1, wantHistory: append(append([]string{}, full[1:]...), "1.gif")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			metrics := &historySamples{}
			store, err := newHistory(t, filepath.Join(t.TempDir(), "history.json"), tc.initial, metrics)
			if err != nil {
				t.Fatal(err)
			}
			source := &retrySource{pages: &retryPages{items: []gallery.Item{{ID: 1}, {ID: 2}}}}
			runner := New(nil, store, source, nil)
			item, err := runner.selectNext(ctx)
			if err != nil || item.ID != tc.wantID {
				t.Fatalf("candidate=%v err=%v", item, err)
			}
			names, err := store.Load(ctx)
			if err != nil || !reflect.DeepEqual(names, tc.initial) {
				t.Fatalf("selection changed history: %v, %v", names, err)
			}
			if len(metrics.sizes) != 0 {
				t.Fatalf("selection emitted metrics: %v", metrics.sizes)
			}
			if err := store.Record(ctx, item.Filename()); err != nil {
				t.Fatal(err)
			}
			names, err = store.Load(ctx)
			if err != nil || !reflect.DeepEqual(names, tc.wantHistory) {
				t.Fatalf("history=%v err=%v", names, err)
			}
			if !reflect.DeepEqual(metrics.sizes, []int{len(tc.wantHistory)}) {
				t.Fatalf("metrics=%v", metrics.sizes)
			}
		})
	}
}

type testDisplay struct {
	failure  error
	locked   bool
	restored bool
	onSend   func()
}

func (d *testDisplay) Lock(context.Context) error { d.locked = true; return nil }
func (d *testDisplay) Unlock()                    { d.locked = false }
func (d *testDisplay) SendAnimation(context.Context, []byte) error {
	if d.onSend != nil {
		d.onSend()
	}
	return d.failure
}
func (d *testDisplay) ShowClock(context.Context, settings.Clock) error {
	d.restored = true
	return nil
}

func TestRunnerLeavesHistoryUnchangedAfterSendFailure(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	metrics := &historySamples{}
	store, err := newHistory(t, filepath.Join(dir, "history.json"), []string{"1.gif"}, metrics)
	if err != nil {
		t.Fatal(err)
	}
	prefs := settings.New(filepath.Join(dir, "settings.json"), settings.Defaults(0, 5*time.Minute, 10*time.Second))
	source := &retrySource{pages: &retryPages{items: []gallery.Item{{ID: 2}}}}
	display := &testDisplay{failure: errors.New("send failed")}
	runner := New(prefs, store, source, display)
	if err := runner.cycle(ctx, 0); !errors.Is(err, display.failure) {
		t.Fatalf("send error lost: %v", err)
	}
	names, err := store.Load(ctx)
	if err != nil || !reflect.DeepEqual(names, []string{"1.gif"}) {
		t.Fatalf("history=%v err=%v", names, err)
	}
	if display.locked || !display.restored {
		t.Fatalf("locked=%v restored=%v", display.locked, display.restored)
	}
	if len(metrics.sizes) != 0 {
		t.Fatalf("metrics=%v", metrics.sizes)
	}
}

type failingHistorySave struct {
	History
	failure error
}

func (h failingHistorySave) Record(context.Context, string) error { return h.failure }

func TestRunnerRecordsOnlyCompletedSends(t *testing.T) {
	for _, scenario := range []string{"success", "cancel after send", "history save failure"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			dir := t.TempDir()
			initial := []string{"1.gif", "2.gif"}
			metrics := &historySamples{}
			store, err := newHistory(t, filepath.Join(dir, "history.json"), initial, metrics)
			if err != nil {
				t.Fatal(err)
			}
			var historyStore History = store
			saveFailure := errors.New("disk full")
			if scenario == "history save failure" {
				historyStore = failingHistorySave{store, saveFailure}
			}
			prefs := settings.New(filepath.Join(dir, "settings.json"), settings.Defaults(0, 5*time.Minute, 10*time.Second))
			source := &retrySource{pages: &retryPages{items: []gallery.Item{{ID: 2}}}}
			display := &testDisplay{onSend: func() {
				names, err := store.Load(ctx)
				if err != nil || !reflect.DeepEqual(names, initial) {
					t.Fatalf("history changed before completed send: %v, %v", names, err)
				}
				if scenario == "cancel after send" {
					cancel()
				}
			}}
			runner := New(prefs, historyStore, source, display)
			duration := time.Duration(0)
			if scenario == "cancel after send" {
				duration = time.Hour
			}
			err = runner.cycle(ctx, duration)
			switch scenario {
			case "success":
				if err != nil {
					t.Fatal(err)
				}
			case "cancel after send":
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("expected cancellation, got %v", err)
				}
			case "history save failure":
				if !errors.Is(err, saveFailure) {
					t.Fatalf("history error lost: %v", err)
				}
			}
			names, err := store.Load(context.Background())
			want := []string{"2.gif"}
			wantMetrics := []int{1}
			if scenario == "history save failure" {
				want = initial
				wantMetrics = nil
			}
			if err != nil || !reflect.DeepEqual(names, want) {
				t.Fatalf("history=%v want=%v err=%v", names, want, err)
			}
			if !reflect.DeepEqual(metrics.sizes, wantMetrics) {
				t.Fatalf("metrics=%v want=%v", metrics.sizes, wantMetrics)
			}
			if display.locked || !display.restored {
				t.Fatalf("locked=%v restored=%v", display.locked, display.restored)
			}
		})
	}
}

func TestRunnerUsesCachedHistoryAcrossCycles(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	metrics := &historySamples{}
	store, err := newHistory(t, path, []string{"1.gif"}, metrics)
	if err != nil {
		t.Fatal(err)
	}
	prefs := settings.New(filepath.Join(dir, "settings.json"), settings.Defaults(0, 5*time.Minute, 10*time.Second))
	source := &retrySource{pages: &retryPages{items: []gallery.Item{{ID: 2}, {ID: 3}}}}
	runner := New(prefs, store, source, &testDisplay{})
	if len(metrics.sizes) != 0 {
		t.Fatalf("startup metrics=%v", metrics.sizes)
	}
	if err := os.WriteFile(path, []byte("invalid external edit"), 0600); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := runner.cycle(ctx, 0); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(source.calls, []int64{2, 3}) {
		t.Fatalf("sent items=%v", source.calls)
	}
	names, err := store.Load(ctx)
	want := []string{"1.gif", "2.gif", "3.gif"}
	if err != nil || !reflect.DeepEqual(names, want) {
		t.Fatalf("persisted history=%v err=%v", names, err)
	}
	if !reflect.DeepEqual(metrics.sizes, []int{2, 3}) {
		t.Fatalf("metrics=%v", metrics.sizes)
	}
}

type testHistory struct {
	*history.Store
	path string
}

func newHistory(t *testing.T, path string, initial []string, metrics history.Metrics) (*testHistory, error) {
	t.Helper()
	if initial != nil {
		data, err := json.Marshal(initial)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	store, err := history.New(path, metrics)
	if err != nil {
		return nil, err
	}
	return &testHistory{Store: store, path: path}, nil
}
func (h *testHistory) Load(context.Context) ([]string, error) {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	err = json.Unmarshal(data, &names)
	return names, err
}
