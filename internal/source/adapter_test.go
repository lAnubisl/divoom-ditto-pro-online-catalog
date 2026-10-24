package source

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/media"
)

type pageRequest struct {
	kind string
	id   int64
	page int
}

type testAPI struct {
	requests  []pageRequest
	pageError error
	items     []gallery.Item
	raw       []byte
	err       error
	downloads int
}

func (a *testAPI) FetchCategoryPage(_ context.Context, category, page int) (gallery.Page, error) {
	a.requests = append(a.requests, pageRequest{"category", int64(category), page})
	return a.page(page)
}
func (a *testAPI) FetchAlbumPage(_ context.Context, albumID int64, page int) (gallery.Page, error) {
	a.requests = append(a.requests, pageRequest{"album", albumID, page})
	return a.page(page)
}
func (a *testAPI) page(page int) (gallery.Page, error) {
	if a.pageError != nil {
		return gallery.Page{}, a.pageError
	}
	if page > 1 {
		return gallery.Page{}, nil
	}
	return gallery.Page{Items: a.items}, nil
}
func (a *testAPI) FetchItem(context.Context, gallery.Item) ([]byte, error) {
	a.downloads++
	return a.raw, a.err
}

func TestFetchPageCachesOnlyLastParameterCombination(t *testing.T) {
	api := &testAPI{items: cursorItems(42)}
	adapter := New(api, nil, nil)
	requests := []pageRequest{
		{"category", 3, 1},
		{"album", 3, 1},
		{"album", 4, 1},
		{"album", 4, 2},
		{"category", 3, 1}, // Returning to an older key fetches it again.
	}
	for i, request := range requests {
		first, err := adapter.fetchPage(context.Background(), request.kind, request.id, request.page)
		if err != nil {
			t.Fatal(err)
		}
		// Change the API response so the second call must return cached data.
		api.items = cursorItems(int64(100 + i))
		cached, err := adapter.fetchPage(context.Background(), request.kind, request.id, request.page)
		if err != nil || !reflect.DeepEqual(cached, first) {
			t.Fatalf("cached page=%+v want=%+v err=%v", cached, first, err)
		}
		if len(api.requests) != i+1 {
			t.Fatalf("unexpected requests: %v", api.requests)
		}
	}
	if !reflect.DeepEqual(api.requests, requests) {
		t.Fatalf("requests=%v want=%v", api.requests, requests)
	}
}

func TestFetchPageRetriesErrors(t *testing.T) {
	api := &testAPI{pageError: errors.New("network timeout")}
	adapter := New(api, nil, nil)
	if _, err := adapter.fetchPage(context.Background(), "category", 3, 1); !errors.Is(err, api.pageError) {
		t.Fatalf("expected page error: %v", err)
	}
	api.pageError = nil
	api.items = cursorItems(42)
	result, err := adapter.fetchPage(context.Background(), "category", 3, 1)
	if err != nil || len(result.Items) != 1 || result.Items[0].ID != 42 || len(api.requests) != 2 {
		t.Fatalf("failed response was not retried: result=%+v requests=%v err=%v", result, api.requests, err)
	}
}

func TestUnsupportedItemRemainsCandidateAfterFailure(t *testing.T) {
	cache, err := media.New(filepath.Join(t.TempDir(), "media"))
	if err != nil {
		t.Fatal(err)
	}
	prefs, _ := cursorStore(t, 1, 0)
	item := gallery.Item{ID: 42, FileID: "artwork", Type: 1, Size: 1}
	api := &testAPI{items: []gallery.Item{item}, raw: []byte{99}}
	adapter := New(api, cache, prefs)
	for attempt := 0; attempt < 2; attempt++ {
		if _, err := adapter.NextAnimation(context.Background()); !errors.Is(err, animation.ErrUnsupportedFormat) {
			t.Fatalf("expected unsupported format: %v", err)
		}
		if got := cursorValues(t, prefs); got.SourcePage != 1 || got.PreviousItemID != item.ID {
			t.Fatalf("failed load did not persist cursor: %+v", got)
		}
		// Recreate the adapter to verify only the cursor survives a restart.
		adapter = New(api, cache, prefs)
	}
	if api.downloads != 2 {
		t.Fatalf("unsupported artwork was not retried on wrap: downloads=%d", api.downloads)
	}
}

func TestNextAnimationSelectionFailureDoesNotDownload(t *testing.T) {
	for _, scenario := range []string{"empty", "page error", "cancelled"} {
		t.Run(scenario, func(t *testing.T) {
			prefs, _ := cursorStore(t, 1, 0)
			cache, err := media.New(filepath.Join(t.TempDir(), "media"))
			if err != nil {
				t.Fatal(err)
			}
			api := &testAPI{}
			want := ErrEmptyGallery
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch scenario {
			case "page error":
				want = errors.New("network timeout")
				api.pageError = want
			case "cancelled":
				want = context.Canceled
				cancel()
			}
			adapter := New(api, cache, prefs)
			if data, err := adapter.NextAnimation(ctx); !errors.Is(err, want) || len(data) != 0 {
				t.Fatalf("unexpected result: bytes=%d err=%v", len(data), err)
			}
			if api.downloads != 0 {
				t.Fatal("selection failure started an artwork download")
			}
		})
	}
}

func TestFailedArtworkIsRetriedOnLaterRequest(t *testing.T) {
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
			cache, err := media.New(filepath.Join(dir, "media"))
			if err != nil {
				t.Fatal(err)
			}
			api := &testAPI{raw: tc.raw, err: tc.err}
			adapter := New(api, cache, nil)
			item := gallery.Item{ID: 42}
			for i := 0; i < 2; i++ {
				if _, err := adapter.Animation(context.Background(), item); err == nil {
					t.Fatal("expected loading error")
				}
			}
			if api.downloads != 2 {
				t.Fatal("failed response was not retried")
			}
		})
	}
}

func TestSupportedArtworkIsDownloadedOnce(t *testing.T) {
	dir := t.TempDir()
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
	item := gallery.Item{ID: 42, FileID: "artwork", Type: 1, Size: 1}
	api := &testAPI{raw: raw, items: []gallery.Item{item}}
	prefs, _ := cursorStore(t, 1, 0)
	adapter := New(api, cache, prefs)
	for i := 0; i < 2; i++ {
		if data, err := adapter.NextAnimation(context.Background()); err != nil || len(data) == 0 {
			t.Fatalf("next animation: bytes=%d err=%v", len(data), err)
		}
	}
	if _, err := adapter.Animation(context.Background(), item); err != nil {
		t.Fatalf("manual animation loading failed: %v", err)
	}
	if api.downloads != 1 {
		t.Fatal("supported artwork must use the existing disk cache")
	}
}
