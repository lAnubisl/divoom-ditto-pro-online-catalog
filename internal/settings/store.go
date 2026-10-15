// Package settings owns the application's shared, persistent preferences.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store is constructed once and shared. Reads load lazily; updates publish only
// after saving successfully, so a disk failure cannot change active preferences.
type Store struct {
	mu     sync.Mutex
	path   string
	values Values
	loaded bool
}

func New(path string, defaults Values) *Store { return &Store{path: path, values: defaults} }

func (s *Store) load() error {
	if s.loaded {
		return nil
	}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.values.validate(); err != nil {
			return err
		}
		if err := s.save(s.values); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		var required struct {
			SourceID   *int64  `json:"source_id"`
			StyleID    *int    `json:"style_id"`
			StyleColor *string `json:"style_color"`
		}
		if err := json.Unmarshal(data, &required); err != nil {
			return err
		}
		if required.SourceID == nil || required.StyleID == nil || required.StyleColor == nil {
			return errors.New("source_id, style_id and style_color are required")
		}
		var v Values
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&v); err != nil {
			return err
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			return errors.New("unexpected data after autoplay settings")
		}
		if err := v.validate(); err != nil {
			return err
		}
		s.values = v
	}
	s.loaded = true
	return nil
}

func (s *Store) Read() (Values, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return Values{}, fmt.Errorf("load settings: %w", err)
	}
	return s.values, nil
}

func (s *Store) Update(change func(*Values)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.load(); err != nil {
		return err
	}
	next := s.values
	change(&next)
	next.StyleColor = strings.ToUpper(next.StyleColor)
	if err := next.validate(); err != nil {
		return err
	}
	if err := s.save(next); err != nil {
		return err
	}
	s.values = next
	return nil
}

func (s *Store) save(v Values) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".settings-*.tmp")
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
