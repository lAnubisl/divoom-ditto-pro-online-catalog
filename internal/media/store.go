package media

import (
	"bytes"
	"context"
	"fmt"
	"image/gif"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// Store retains valid GIFs indefinitely. Downloading and conversion belong to
// the source adapter, not disk storage.
type Store struct{ dir string }

func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}
func (s *Store) path(id int64) string { return filepath.Join(s.dir, strconv.FormatInt(id, 10)+".gif") }
func (s *Store) Load(ctx context.Context, id int64) ([]byte, bool) {
	if ctx.Err() != nil || id <= 0 {
		return nil, false
	}
	f, err := os.Open(s.path(id))
	if err != nil {
		return nil, false
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, MaxAnimationSize+1))
	if err != nil || len(data) > MaxAnimationSize {
		return nil, false
	}
	decoded, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil || decoded.Config.Width != 16 || decoded.Config.Height != 16 || len(decoded.Image) == 0 {
		return nil, false
	}
	return data, true
}
func (s *Store) Save(ctx context.Context, id int64, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("invalid animation ID")
	}
	f, err := os.CreateTemp(s.dir, ".animation-*.tmp")
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
	return os.Rename(f.Name(), s.path(id))
}
