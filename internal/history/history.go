// Package history persists the names of automatically displayed animations.
package history

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct{ path string }

func New(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	h := &Store{path: path}
	_, err := h.Load(context.Background())
	return h, err
}
func (h *Store) Load(ctx context.Context) ([]string, error) {
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
func (h *Store) Save(ctx context.Context, names []string) error {
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
