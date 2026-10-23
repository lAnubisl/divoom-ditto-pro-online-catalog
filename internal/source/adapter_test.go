package source

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/media"
	"divoom-gallery-client/internal/settings"
	"divoom-gallery-client/internal/skipped"
)

type testAPI struct {
	items     []gallery.Item
	raw       []byte
	err       error
	downloads int
}

func (a *testAPI) FetchCategoryPage(_ context.Context, _, page int) (gallery.Page, error) {
	if page > 1 {
		return gallery.Page{}, nil
	}
	return gallery.Page{Items: a.items}, nil
}
func (a *testAPI) FetchAlbumPage(ctx context.Context, _ int64, page int) (gallery.Page, error) {
	return a.FetchCategoryPage(ctx, 0, page)
}
func (a *testAPI) FetchItem(context.Context, gallery.Item) ([]byte, error) {
	a.downloads++
	return a.raw, a.err
}

func TestUnsupportedItemsPersistWithoutBreakingPagination(t *testing.T) {
	dir := t.TempDir()
	prefs := settings.New(filepath.Join(dir, "settings.json"), settings.Defaults(0, time.Minute, time.Second))
	cache, err := media.New(filepath.Join(dir, "media"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "skipped.json")
	store, err := skipped.New(path)
	if err != nil {
		t.Fatal(err)
	}
	item := gallery.Item{ID: 42, FileID: "artwork", Type: 1, Size: 1}
	api := &testAPI{items: []gallery.Item{item}, raw: []byte{99}}
	adapter := New(prefs, api, cache, store)
	ctx := context.Background()
	if _, err := adapter.Animation(ctx, item); !errors.Is(err, animation.ErrUnsupportedFormat) {
		t.Fatalf("expected unsupported encoding, got %v", err)
	}
	store, err = skipped.New(path)
	if err != nil {
		t.Fatal(err)
	}
	adapter = New(prefs, api, cache, store)
	if _, err := adapter.Animation(ctx, item); !errors.Is(err, animation.ErrUnsupportedFormat) || api.downloads != 1 {
		t.Fatalf("known unsupported item downloaded again: downloads=%d err=%v", api.downloads, err)
	}
	for _, source := range []string{"category", "album"} {
		if err := prefs.Update(func(v *settings.Values) { v.Source = source; v.SourceID = 34 }); err != nil {
			t.Fatal(err)
		}
		pages, err := adapter.Pages()
		if err != nil {
			t.Fatal(err)
		}
		items, end, err := pages(ctx, 1)
		if err != nil || len(items) != 0 || end {
			t.Fatalf("filtered page must not end pagination: items=%v end=%v err=%v", items, end, err)
		}
		_, end, err = pages(ctx, 2)
		if err != nil || !end {
			t.Fatalf("empty remote page must end pagination: end=%v err=%v", end, err)
		}
	}
}

func TestTransientAndMalformedArtworkAreNotSkipped(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  []byte
		err  error
	}{
		{name: "network", err: errors.New("network timeout")},
		{name: "empty"},
		{name: "truncated known format", raw: []byte{9}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			store, err := skipped.New(filepath.Join(dir, "skipped.json"))
			if err != nil {
				t.Fatal(err)
			}
			cache, err := media.New(filepath.Join(dir, "media"))
			if err != nil {
				t.Fatal(err)
			}
			api := &testAPI{raw: tc.raw, err: tc.err}
			adapter := New(nil, api, cache, store)
			item := gallery.Item{ID: 42}
			for i := 0; i < 2; i++ {
				if _, err := adapter.Animation(context.Background(), item); err == nil {
					t.Fatal("expected loading error")
				}
			}
			if store.Contains(item.ID) || api.downloads != 2 {
				t.Fatal("temporary or malformed response was permanently skipped")
			}
		})
	}
}

func TestSupportedArtworkIsDownloadedOnce(t *testing.T) {
	dir := t.TempDir()
	store, err := skipped.New(filepath.Join(dir, "skipped.json"))
	if err != nil {
		t.Fatal(err)
	}
	cache, err := media.New(filepath.Join(dir, "media"))
	if err != nil {
		t.Fatal(err)
	}
	raw := make([]byte, 769)
	raw[0] = 8
	block, err := aes.NewCipher([]byte("78hrey23y28ogs89"))
	if err != nil {
		t.Fatal(err)
	}
	cipher.NewCBCEncrypter(block, []byte("1234567890123456")).CryptBlocks(raw[1:], make([]byte, 768))
	api := &testAPI{raw: raw}
	adapter := New(nil, api, cache, store)
	for i := 0; i < 2; i++ {
		if _, err := adapter.Animation(context.Background(), gallery.Item{ID: 42}); err != nil {
			t.Fatal(err)
		}
	}
	if api.downloads != 1 || store.Contains(42) {
		t.Fatal("supported artwork must use the existing disk cache")
	}
}
