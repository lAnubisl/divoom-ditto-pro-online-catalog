package skipped

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentAddsPersistUniqueIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skipped.json")
	store, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			if err := store.Add(context.Background(), id); err != nil {
				t.Error(err)
			}
		}(int64(i%5 + 1))
	}
	wg.Wait()
	reloaded, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	for id := int64(1); id <= 5; id++ {
		if !reloaded.Contains(id) {
			t.Fatalf("missing ID %d after restart", id)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var ids []int64
	if err := json.Unmarshal(data, &ids); err != nil || len(ids) != 5 {
		t.Fatalf("expected five unique IDs: %s, err=%v", data, err)
	}
}

func TestFailedSaveDoesNotPublishID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skipped.json")
	store, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := store.Add(context.Background(), 42); err == nil || store.Contains(42) {
		t.Fatal("failed write must not change the active skip list")
	}
}
