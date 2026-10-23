// Package skipped persists gallery IDs whose artwork encoding is unsupported.
package skipped

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store is shared by previews and rotation. Updates are serialized and become
// visible only after the new JSON array has been saved successfully.
type Store struct {
	mu    sync.Mutex
	path  string
	ids   []int64
	known map[int64]bool
}

func New(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	s := &Store{path: path, ids: []int64{}, known: make(map[int64]bool)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.save(s.ids); err != nil {
			return nil, err
		}
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.ids); err != nil {
		return nil, fmt.Errorf("invalid skipped items: %w", err)
	}
	for _, id := range s.ids {
		if id <= 0 {
			return nil, errors.New("invalid skipped gallery ID")
		}
		s.known[id] = true
	}
	return s, nil
}

func (s *Store) Contains(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.known[id]
}

func (s *Store) Add(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("invalid skipped gallery ID")
	}
	if s.known[id] {
		return nil
	}
	next := append(append([]int64{}, s.ids...), id)
	if err := s.save(next); err != nil {
		return err
	}
	s.ids, s.known[id] = next, true
	return nil
}

func (s *Store) save(ids []int64) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".skipped-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(append(data, '\n')); err != nil {
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
	return os.Rename(f.Name(), s.path)
}
