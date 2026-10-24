package rotation

import (
	"context"
	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/media"
	"divoom-gallery-client/internal/settings"
	"divoom-gallery-client/internal/source"
	"errors"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func pageItems(ids ...int64) []gallery.Item {
	items := make([]gallery.Item, len(ids))
	for i, id := range ids {
		items[i] = gallery.Item{ID: id, FileID: "artwork", Type: 1, Size: 1}
	}
	return items
}
func readProgress(t *testing.T, prefs *settings.Store) settings.Values {
	t.Helper()
	v, err := prefs.Read()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

type unsupportedAPI struct {
	pages     []int
	downloads []int64
}

func (a *unsupportedAPI) FetchCategoryPage(_ context.Context, _, page int) (gallery.Page, error) {
	a.pages = append(a.pages, page)
	if page > 1 {
		return gallery.Page{}, nil
	}
	return gallery.Page{Items: pageItems(11, 12)}, nil
}
func (a *unsupportedAPI) FetchAlbumPage(ctx context.Context, _ int64, page int) (gallery.Page, error) {
	return a.FetchCategoryPage(ctx, 0, page)
}
func (a *unsupportedAPI) FetchItem(_ context.Context, item gallery.Item) ([]byte, error) {
	a.downloads = append(a.downloads, item.ID)
	return []byte{99}, nil
}

func TestRealAdapterAdvancesAfterUnsupportedAnimationAndRetriesOnWrap(t *testing.T) {
	prefs := newPreferences(t)
	cache, err := media.New(filepath.Join(t.TempDir(), "media"))
	if err != nil {
		t.Fatal(err)
	}
	api := &unsupportedAPI{}
	adapter := source.New(api, cache, prefs)
	display := &testDisplay{}
	runner := New(prefs, adapter, display)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i, id := range []int64{11, 12, 11} {
		if err := runner.cycle(ctx, 0); !errors.Is(err, animation.ErrUnsupportedFormat) {
			t.Fatalf("expected unsupported candidate: %v", err)
		}
		if len(api.downloads) != i+1 || api.downloads[i] != id {
			t.Fatalf("failed iteration did not advance one candidate: %v", api.downloads)
		}
		progress := readProgress(t, prefs)
		if progress.SourcePage != 1 || progress.PreviousItemID != id {
			t.Fatalf("failed candidate cursor was not retained: %+v", progress)
		}
		if display.locked || display.sends != 0 {
			t.Fatalf("failed load retained lock or uploaded artwork: %+v", display)
		}
		// Recreate the adapter; failures do not exclude the item on a later wrap.
		runner = New(prefs, source.New(api, cache, prefs), display)
	}
	if !slices.Equal(api.pages, []int{1, 1, 1, 2, 1}) {
		t.Fatalf("unexpected pages: %v", api.pages)
	}
}
