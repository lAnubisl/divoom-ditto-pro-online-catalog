// Package history persists the names of automatically displayed animations.
package history

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

const Limit = 1000

// Metrics receives the retained history size after a successful save.
type Metrics interface {
	WriteHistorySize(context.Context, int) error
}

// Store owns playback history in memory and on disk. External file edits are
// picked up only at startup; failed saves leave the active history unchanged.
type Store struct {
	mu      sync.Mutex
	path    string
	played  []string
	metrics Metrics
}

func New(path string, metrics Metrics) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	h := &Store{path: path, metrics: metrics}
	played, err := h.load(context.Background())
	if err != nil {
		return nil, fmt.Errorf("load playback history: %w", err)
	}
	h.played = played
	return h, nil
}

func (h *Store) Contains(name string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return slices.Contains(h.played, name)
}

// Record remembers a completed send. Recording an already played name restarts
// an exhausted source. Metrics failures do not invalidate successful persistence.
func (h *Store) Record(ctx context.Context, name string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	played := slices.Clone(h.played)
	if slices.Contains(played, name) {
		played = nil
	}
	played = append(played, name)
	if len(played) > Limit {
		played = played[len(played)-Limit:]
	}
	if err := h.save(ctx, played); err != nil {
		return err
	}
	h.played = played
	if h.metrics != nil {
		if err := h.metrics.WriteHistorySize(ctx, len(played)); err != nil {
			log.Printf("report history size: %v", err)
		}
	}
	return nil
}

func (h *Store) load(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	if err := json.Unmarshal(data, &names); err != nil {
		return nil, fmt.Errorf("invalid playback history: %w", err)
	}
	return names, nil
}
func (h *Store) save(ctx context.Context, names []string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(append([]string{}, names...))
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(h.path), ".history-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), h.path)
}
