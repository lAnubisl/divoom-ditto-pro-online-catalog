package settings_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"divoom-gallery-client/internal/settings"
)

func defaults() settings.Values {
	return settings.Defaults(0, time.Minute, time.Second)
}

func writeValues(t *testing.T, path string, v settings.Values) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func newStore(t *testing.T, path string, initial settings.Values) *settings.Store {
	t.Helper()
	store, err := settings.New(path, initial)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func readValues(t *testing.T, store *settings.Store) settings.Values {
	t.Helper()
	v, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestNewLoadsAndReadKeepsSnapshot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	want := defaults()
	want.SourceID = 42
	writeValues(t, path, want)
	store := newStore(t, path, defaults())
	// Corrupt the file before even the first Read: New must have loaded it.
	if err := os.WriteFile(path, []byte("invalid JSON"), 0600); err != nil {
		t.Fatal(err)
	}
	got := readValues(t, store)
	if got != want {
		t.Fatalf("first read = %+v, want %+v", got, want)
	}
	got.SourceName = "caller mutation"
	if got := readValues(t, store); got != want {
		t.Fatalf("cached read = %+v, want %+v", got, want)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if got := readValues(t, store); got != want {
		t.Fatalf("read after deletion = %+v", got)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("cached read recreated file: %v", err)
	}
}

func TestUpdatePersistsSnapshot(t *testing.T) {
	for _, existing := range []bool{false, true} {
		name := "missing file"
		if existing {
			name = "existing file"
		}
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "settings.json")
			want := defaults()
			if existing {
				want.SourceID = 42
				writeValues(t, path, want)
			}
			store := newStore(t, path, defaults())
			var retained *settings.Values
			if err := store.Update(func(v *settings.Values) {
				v.SourceName = "updated"
				v.StyleColor = "#abcdef"
				retained = v
			}); err != nil {
				t.Fatal(err)
			}
			want.SourceName = "updated"
			want.StyleColor = "#ABCDEF"
			retained.SourceName = "later callback mutation"
			if got := readValues(t, store); got != want {
				t.Fatalf("memory = %+v, want %+v", got, want)
			}
			if got := readValues(t, newStore(t, path, defaults())); got != want {
				t.Fatalf("disk = %+v, want %+v", got, want)
			}
		})
	}
}

func TestNewPersistsDefaultsForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "settings.json")
	store := newStore(t, path, defaults())
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("constructor did not persist defaults: %v", err)
	}
	if got := readValues(t, store); got != defaults() {
		t.Fatalf("defaults = %+v", got)
	}
	if got := readValues(t, newStore(t, path, defaults())); got != defaults() {
		t.Fatalf("persisted defaults = %+v", got)
	}
}

func TestNewRejectsInvalidSettings(t *testing.T) {
	valid, err := json.Marshal(defaults())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		data string
	}{
		{name: "malformed", data: "invalid JSON"},
		{name: "missing required fields", data: `{ "version": 2 }`},
		{name: "unknown field", data: string(valid[:len(valid)-1]) + `,"unknown":true}`},
		{name: "invalid values", data: string(valid[:len(valid)-1]) + `,"duration":"invalid"}`},
		{name: "trailing data", data: string(valid) + `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "settings.json")
			if err := os.WriteFile(path, []byte(tc.data), 0600); err != nil {
				t.Fatal(err)
			}
			store, err := settings.New(path, defaults())
			if err == nil || store != nil {
				t.Fatalf("invalid settings accepted: store=%v err=%v", store, err)
			}
			data, err := os.ReadFile(path)
			if err != nil || string(data) != tc.data {
				t.Fatalf("constructor changed invalid file: data=%q err=%v", data, err)
			}
			want := defaults()
			want.SourceID = 7
			writeValues(t, path, want)
			if got := readValues(t, newStore(t, path, defaults())); got != want {
				t.Fatalf("new constructor after repair = %+v, want %+v", got, want)
			}
		})
	}
}

func TestNewRejectsReadFailure(t *testing.T) {
	store, err := settings.New(t.TempDir(), defaults())
	if err == nil || store != nil {
		t.Fatalf("directory accepted as settings: store=%v err=%v", store, err)
	}
}

func TestNewRejectsInvalidDefaultsForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	invalid := defaults()
	invalid.Duration = "invalid"
	store, err := settings.New(path, invalid)
	if err == nil || store != nil {
		t.Fatalf("invalid defaults accepted: store=%v err=%v", store, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("invalid defaults persisted: %v", err)
	}
}

func TestNewReportsInitialSaveFailure(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "blocked")
	if err := os.WriteFile(parent, []byte("file blocks directory creation"), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(parent, "settings.json")
	store, err := settings.New(path, defaults())
	if err == nil || store != nil {
		t.Fatalf("expected constructor failure: store=%v err=%v", store, err)
	}
	if err := os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	if got := readValues(t, newStore(t, path, defaults())); got != defaults() {
		t.Fatalf("new constructor after repair = %+v", got)
	}
	if got := readValues(t, newStore(t, path, defaults())); got != defaults() {
		t.Fatalf("persisted defaults = %+v", got)
	}
}

func TestRejectedUpdatePreservesSnapshot(t *testing.T) {
	for _, failure := range []string{"validation", "save"} {
		t.Run(failure, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "settings.json")
			store := newStore(t, path, defaults())
			before := readValues(t, store)
			if failure == "save" {
				// A directory cannot be replaced by the atomic file rename.
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			}
			err := store.Update(func(v *settings.Values) {
				v.SourceName = "must not be published"
				if failure == "validation" {
					v.Duration = "invalid"
				}
			})
			if err == nil {
				t.Fatal("expected rejected update")
			}
			if got := readValues(t, store); got != before {
				t.Fatalf("failed update changed memory: %+v", got)
			}
			if failure == "save" {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else if got := readValues(t, newStore(t, path, defaults())); got != before {
				t.Fatalf("failed validation changed disk: %+v", got)
			}
			if err := store.Update(func(v *settings.Values) { v.SourceName = "retry" }); err != nil {
				t.Fatal(err)
			}
			before.SourceName = "retry"
			if got := readValues(t, store); got != before {
				t.Fatalf("retry memory = %+v", got)
			}
			if got := readValues(t, newStore(t, path, defaults())); got != before {
				t.Fatalf("retry disk = %+v", got)
			}
		})
	}
}

func TestConcurrentReadsAndUpdates(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	initial := defaults()
	initial.Source = "album"
	initial.SourceID = 1
	store := newStore(t, path, initial)
	const workers, updates = 8, 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < updates; j++ {
				v, err := store.Read()
				if err != nil {
					t.Error(err)
					return
				}
				if v.Source != "album" || v.SourceID < 1 {
					t.Errorf("invalid snapshot: %+v", v)
					return
				}
				if err := store.Update(func(v *settings.Values) { v.SourceID++ }); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	close(start)
	wg.Wait()
	want := initial
	want.SourceID += workers * updates
	if got := readValues(t, store); got != want {
		t.Fatalf("lost updates: got %+v, want %+v", got, want)
	}
	if got := readValues(t, newStore(t, path, defaults())); got != want {
		t.Fatalf("disk differs: got %+v, want %+v", got, want)
	}
}

func TestSourcePageLegacyJSONAndRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	legacy := `{"version":2,"enabled":true,"source":"category","source_id":3,"interval":"1m0s","duration":"1s","style_id":0,"style_color":"#00FF00"}`
	if err := os.WriteFile(path, []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	store := newStore(t, path, defaults())
	before := readValues(t, store)
	if before.SourcePage != 0 || before.PreviousItemID != 0 {
		t.Fatalf("legacy page=%d", before.SourcePage)
	}
	if err := store.Update(func(v *settings.Values) { v.SourcePage = 7; v.PreviousItemID = 71 }); err != nil {
		t.Fatal(err)
	}
	want := before
	want.SourcePage = 7
	want.PreviousItemID = 71
	if got := readValues(t, newStore(t, path, defaults())); got != want {
		t.Fatalf("restart=%+v, want %+v", got, want)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["source_page"]) != "7" || string(fields["previous_item_id"]) != "71" {
		t.Fatalf("page not persisted: %s", data)
	}
	if err := store.Update(func(v *settings.Values) { v.SourcePage = -1 }); err == nil {
		t.Fatal("negative page accepted")
	}
	if err := store.Update(func(v *settings.Values) { v.PreviousItemID = -1 }); err == nil {
		t.Fatal("negative previous item ID accepted")
	}
	if got := readValues(t, store); got != want {
		t.Fatalf("invalid update changed memory: %+v", got)
	}
	if got := readValues(t, newStore(t, path, defaults())); got != want {
		t.Fatalf("invalid update changed disk: %+v", got)
	}
}
