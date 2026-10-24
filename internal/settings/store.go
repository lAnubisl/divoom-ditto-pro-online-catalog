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

// Store is constructed once and shared. New loads settings; updates publish only
// after saving successfully, so a disk failure cannot change active preferences.
type Store struct {
	mu     sync.Mutex
	path   string
	values Values
}

// New loads and validates settings, creating the file from defaults if missing.
func New(path string, defaults Values) (*Store, error) {
	s := &Store{path: path, values: defaults}
	if err := s.load(); err != nil {
		return nil, fmt.Errorf("load settings: %w", err)
	}
	return s, nil
}

// load initializes the store before New exposes it to callers.
func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s.saveDefaults()
	}
	if err != nil {
		return err
	}
	v, err := decodeValues(data)
	if err != nil {
		return err
	}
	s.values = v
	return nil
}

func (s *Store) saveDefaults() error {
	if err := s.values.validate(); err != nil {
		return err
	}
	return s.save(s.values)
}

func decodeValues(data []byte) (Values, error) {
	if err := validateRequiredFields(data); err != nil {
		return Values{}, err
	}
	var v Values
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&v); err != nil {
		return Values{}, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return Values{}, errors.New("unexpected data after autoplay settings")
	}
	if err := v.validate(); err != nil {
		return Values{}, err
	}
	return v, nil
}

func validateRequiredFields(data []byte) error {
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
	return nil
}

// Read returns a copy of the in-memory values.
// Values contains no mutable reference fields.
func (s *Store) Read() (Values, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values, nil
}

// Update serializes changes with reads and other updates. Validation or save
// failures leave the active values unchanged.
func (s *Store) Update(change func(*Values)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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
