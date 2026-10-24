package source

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/settings"
)

type cursorAPI struct {
	testAPI
	pages  map[int][]gallery.Item
	onPage func()
}

func (a *cursorAPI) FetchCategoryPage(_ context.Context, id, page int) (gallery.Page, error) {
	return a.fetch("category", int64(id), page)
}
func (a *cursorAPI) FetchAlbumPage(_ context.Context, id int64, page int) (gallery.Page, error) {
	return a.fetch("album", id, page)
}
func (a *cursorAPI) fetch(kind string, id int64, page int) (gallery.Page, error) {
	a.requests = append(a.requests, pageRequest{kind, id, page})
	if a.onPage != nil {
		a.onPage()
	}
	return gallery.Page{Items: a.pages[page]}, a.pageError
}

func cursorItems(ids ...int64) []gallery.Item {
	items := make([]gallery.Item, len(ids))
	for i, id := range ids {
		items[i] = gallery.Item{ID: id, FileID: "artwork", Type: 1, Size: 1}
	}
	return items
}

func cursorStore(t *testing.T, page int, previous int64) (*settings.Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "settings.json")
	v := settings.Defaults(3, time.Minute, time.Second)
	v.SourcePage, v.PreviousItemID, v.SourceName = page, previous, "keep name"
	prefs, err := settings.New(path, v)
	if err != nil {
		t.Fatal(err)
	}
	return prefs, path
}

func cursorAdapter(t *testing.T, api API, prefs *settings.Store) *Adapter {
	t.Helper()
	return New(api, nil, prefs)
}

func cursorValues(t *testing.T, prefs *settings.Store) settings.Values {
	t.Helper()
	v, err := prefs.Read()
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func wantCandidate(t *testing.T, a *Adapter, id int64, page int) {
	t.Helper()
	item, err := a.nextCandidate(context.Background())
	if err != nil || item.ID != id {
		t.Fatalf("candidate=%v want=%d err=%v", item, id, err)
	}
	v := cursorValues(t, a.settings)
	if v.SourcePage != page || v.PreviousItemID != id {
		t.Fatalf("cursor not persisted: %+v", v)
	}
}

func TestCandidateResumesCursorAcrossRestartAndWraps(t *testing.T) {
	prefs, path := cursorStore(t, 3, 31)
	before := cursorValues(t, prefs)
	api := &cursorAPI{pages: map[int][]gallery.Item{1: cursorItems(11), 3: cursorItems(31, 32), 4: cursorItems(41)}}
	a := cursorAdapter(t, api, prefs)
	wantCandidate(t, a, 32, 3)
	reopened, err := settings.New(path, settings.Defaults(0, time.Minute, time.Second))
	if err != nil {
		t.Fatal(err)
	}
	a = New(api, nil, reopened)
	wantCandidate(t, a, 41, 4)
	wantCandidate(t, a, 11, 1)
	want := []pageRequest{{"category", 3, 3}, {"category", 3, 3}, {"category", 3, 4}, {"category", 3, 5}, {"category", 3, 1}}
	if !slices.Equal(api.requests, want) {
		t.Fatalf("pages=%v want=%v", api.requests, want)
	}
	after := cursorValues(t, reopened)
	after.SourcePage, after.PreviousItemID = before.SourcePage, before.PreviousItemID
	if after != before {
		t.Fatalf("unrelated settings changed: %+v", after)
	}
}

func TestSourceEndReturnsFirstCandidateWithoutError(t *testing.T) {
	for _, kind := range []string{"category", "album"} {
		for _, startPage := range []int{2, 9} {
			t.Run(kind+"/"+strconv.Itoa(startPage), func(t *testing.T) {
				prefs, _ := cursorStore(t, startPage, 22)
				if err := prefs.Update(func(v *settings.Values) { v.Source = kind }); err != nil {
					t.Fatal(err)
				}
				api := &cursorAPI{pages: map[int][]gallery.Item{
					1: cursorItems(11, 12),
					2: cursorItems(21, 22),
				}}
				a := cursorAdapter(t, api, prefs)
				// Both an exhausted final page and a stale cursor beyond the end
				// return the first candidate directly, including on later wraps.
				for pass := 0; pass < 3; pass++ {
					wantCandidate(t, a, 11, 1)
					wantCandidate(t, a, 12, 1)
					wantCandidate(t, a, 21, 2)
					wantCandidate(t, a, 22, 2)
				}
			})
		}
	}
}

func TestCandidateFiltersWithoutLosingRawCursorOrPageBoundary(t *testing.T) {
	for _, kind := range []string{"category", "album"} {
		t.Run(kind, func(t *testing.T) {
			prefs, _ := cursorStore(t, 1, 1)
			if err := prefs.Update(func(v *settings.Values) { v.Source = kind }); err != nil {
				t.Fatal(err)
			}
			api := &cursorAPI{pages: map[int][]gallery.Item{1: cursorItems(1), 2: cursorItems(2, 3), 3: cursorItems(4)}}
			api.pages[2][0].Type = 0 // Album metadata allows type 0.
			api.pages[2][1].Size = 2
			a := cursorAdapter(t, api, prefs)
			if kind == "album" {
				wantCandidate(t, a, 2, 2)
			} else {
				wantCandidate(t, a, 4, 3)
			}
			for _, req := range api.requests {
				if req.kind != kind || req.id != 3 {
					t.Fatalf("wrong source: %v", req)
				}
			}
		})
	}
}

func TestMissingCursorRestartsCurrentPage(t *testing.T) {
	prefs, _ := cursorStore(t, 3, 999)
	a := cursorAdapter(t, &cursorAPI{pages: map[int][]gallery.Item{3: cursorItems(31, 32)}}, prefs)
	wantCandidate(t, a, 31, 3)
}

func TestEmptyAndFilteredSourcesResetCursor(t *testing.T) {
	for _, filtered := range []bool{false, true} {
		prefs, _ := cursorStore(t, 3, 31)
		api := &cursorAPI{pages: map[int][]gallery.Item{}}
		a := cursorAdapter(t, api, prefs)
		if filtered {
			api.pages[1], api.pages[3] = cursorItems(11), cursorItems(31)
			api.pages[1][0].Size, api.pages[3][0].Size = 2, 2
		}
		if _, err := a.nextCandidate(context.Background()); !errors.Is(err, ErrEmptyGallery) {
			t.Fatalf("expected empty: %v", err)
		}
		v := cursorValues(t, prefs)
		if v.SourcePage != 0 || v.PreviousItemID != 0 {
			t.Fatalf("cursor not reset: %+v", v)
		}
	}
}

func TestCandidateFailuresPreserveCursor(t *testing.T) {
	for _, failure := range []string{"page", "cancel", "save", "limit"} {
		t.Run(failure, func(t *testing.T) {
			page := 3
			if failure == "limit" {
				page = 10000
			}
			prefs, path := cursorStore(t, page, 31)
			before := cursorValues(t, prefs)
			api := &cursorAPI{pages: map[int][]gallery.Item{page: cursorItems(31, 32)}}
			a := cursorAdapter(t, api, prefs)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch failure {
			case "page":
				api.pageError = errors.New("network timeout")
			case "cancel":
				api.onPage = cancel
			case "save":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			case "limit":
				api.pages[page] = cursorItems(31)
			}
			_, err := a.nextCandidate(ctx)
			if err == nil || errors.Is(err, ErrEmptyGallery) {
				t.Fatalf("expected failure: %v", err)
			}
			if failure == "page" && !errors.Is(err, api.pageError) {
				t.Fatalf("lost page error: %v", err)
			}
			if failure == "cancel" && !errors.Is(err, context.Canceled) {
				t.Fatalf("lost cancellation: %v", err)
			}
			if failure == "save" && !strings.Contains(err.Error(), "save source cursor") {
				t.Fatalf("lost save error: %v", err)
			}
			if got := cursorValues(t, prefs); got != before {
				t.Fatalf("failed request changed cursor: %+v", got)
			}
		})
	}
}

func TestOldSourceCannotOverwriteNewCursor(t *testing.T) {
	for _, atEnd := range []bool{false, true} {
		prefs, _ := cursorStore(t, 3, 0)
		api := &cursorAPI{pages: map[int][]gallery.Item{3: cursorItems(31)}}
		if atEnd {
			api.pages = nil
		}
		a := cursorAdapter(t, api, prefs)
		api.onPage = func() {
			if err := prefs.Update(func(v *settings.Values) { v.Source = "album"; v.SourceID = 7; v.SourcePage = 9; v.PreviousItemID = 91 }); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := a.nextCandidate(context.Background()); err == nil {
			t.Fatal("accepted old source")
		}
		v := cursorValues(t, prefs)
		if v.Source != "album" || v.SourceID != 7 || v.SourcePage != 9 || v.PreviousItemID != 91 {
			t.Fatalf("overwrote new source: %+v", v)
		}
		if len(api.requests) != 1 {
			t.Fatalf("continued stale traversal: %v", api.requests)
		}
	}
}
