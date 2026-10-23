package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type samples struct {
	sizes   []int
	failure error
}

func (m *samples) WriteHistorySize(_ context.Context, size int) error {
	m.sizes = append(m.sizes, size)
	return m.failure
}

func TestStartup(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		wantError  bool
	}{
		{name: "missing"},
		{name: "malformed", data: "invalid JSON", wantError: true},
		{name: "existing", data: `["1.gif"]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "history.json")
			if tc.data != "" {
				if err := os.WriteFile(path, []byte(tc.data), 0600); err != nil {
					t.Fatal(err)
				}
			}
			metrics := &samples{}
			store, err := New(path, metrics)
			if tc.wantError {
				if err == nil || store != nil {
					t.Fatalf("malformed history accepted: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if store.Contains("1.gif") != (tc.name == "existing") {
				t.Fatal("incorrect loaded history")
			}
			if len(metrics.sizes) != 0 {
				t.Fatalf("startup emitted metrics: %v", metrics.sizes)
			}
		})
	}
}

func TestRecordPersistsAndReloads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "history.json")
	metrics := &samples{}
	store, err := New(path, metrics)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"1.gif", "2.gif", "1.gif"} {
		if err := store.Record(context.Background(), name); err != nil {
			t.Fatal(err)
		}
	}
	reopened, err := New(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !store.Contains("1.gif") || store.Contains("2.gif") || !reopened.Contains("1.gif") || reopened.Contains("2.gif") {
		t.Fatal("replayed item did not reset active and persisted history")
	}
	if !reflect.DeepEqual(metrics.sizes, []int{1, 2, 1}) {
		t.Fatalf("metrics=%v", metrics.sizes)
	}
}

func TestRecordFailurePreservesHistory(t *testing.T) {
	for _, scenario := range []string{"cancelled", "persistence failure"} {
		t.Run(scenario, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "history.json")
			metrics := &samples{}
			store, err := New(path, metrics)
			if err != nil {
				t.Fatal(err)
			}
			if err := store.Record(context.Background(), "1.gif"); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if scenario == "cancelled" {
				cancel()
			} else {
				// A directory cannot be replaced by the atomic file rename.
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			}
			if err := store.Record(ctx, "2.gif"); err == nil {
				t.Fatal("expected record failure")
			}
			if !store.Contains("1.gif") || store.Contains("2.gif") {
				t.Fatal("failed save changed active history")
			}
			if !reflect.DeepEqual(metrics.sizes, []int{1}) {
				t.Fatalf("failed save emitted metrics: %v", metrics.sizes)
			}
			if scenario == "cancelled" {
				reopened, err := New(path, nil)
				if err != nil {
					t.Fatal(err)
				}
				if !reopened.Contains("1.gif") || reopened.Contains("2.gif") {
					t.Fatal("cancelled save changed persisted history")
				}
			}
		})
	}
}

func TestMetricsFailureDoesNotUndoRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	metrics := &samples{failure: errors.New("metrics unavailable")}
	store, err := New(path, metrics)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(context.Background(), "1.gif"); err != nil {
		t.Fatal(err)
	}
	reopened, err := New(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !store.Contains("1.gif") || !reopened.Contains("1.gif") {
		t.Fatal("metrics failure lost successful record")
	}
}

func TestRetentionLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	initial := make([]string, Limit)
	for i := range initial {
		initial[i] = fmt.Sprintf("%d.gif", i)
	}
	data, err := json.Marshal(initial)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	metrics := &samples{}
	store, err := New(path, metrics)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(context.Background(), "new.gif"); err != nil {
		t.Fatal(err)
	}
	reopened, err := New(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if store.Contains(initial[0]) || reopened.Contains(initial[0]) {
		t.Fatal("oldest item was not evicted")
	}
	if !reopened.Contains(initial[1]) || !reopened.Contains("new.gif") {
		t.Fatal("retained entries were lost")
	}
	if !reflect.DeepEqual(metrics.sizes, []int{Limit}) {
		t.Fatalf("metrics=%v", metrics.sizes)
	}
}
