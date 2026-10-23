package rotation

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/history"
)

type retryPages struct{ items []gallery.Item }

func (p retryPages) Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error) {
	return func(_ context.Context, page int) ([]gallery.Item, bool, error) {
		if page > 1 {
			return nil, true, nil
		}
		return p.items, true, nil
	}, nil
}

// Models the source's persistent filtering without requiring device access.
type retrySource struct {
	pages   *retryPages
	failure error
	calls   []int64
}

func (s *retrySource) Animation(_ context.Context, item gallery.Item) ([]byte, error) {
	s.calls = append(s.calls, item.ID)
	if item.ID == 1 {
		if errors.Is(s.failure, animation.ErrUnsupportedFormat) {
			s.pages.items = s.pages.items[1:]
		}
		return nil, s.failure
	}
	return []byte("cached GIF"), nil
}

func TestRotationImmediatelyAdvancesPastUnsupportedItem(t *testing.T) {
	for _, allUnsupported := range []bool{false, true} {
		t.Run(map[bool]string{false: "next supported", true: "all unsupported"}[allUnsupported], func(t *testing.T) {
			store, err := history.New(filepath.Join(t.TempDir(), "history.json"))
			if err != nil {
				t.Fatal(err)
			}
			pages := &retryPages{items: []gallery.Item{{ID: 1}, {ID: 2}}}
			if allUnsupported {
				pages.items = pages.items[:1]
			}
			source := &retrySource{pages: pages, failure: animation.ErrUnsupportedFormat}
			runner := &Runner{selection: NewSelector(store, pages, nil), source: source}
			item, data, err := runner.nextAnimation(context.Background())
			if allUnsupported {
				if !errors.Is(err, ErrEmptyGallery) {
					t.Fatalf("expected empty gallery, got %v", err)
				}
			} else if err != nil || item.ID != 2 || len(data) == 0 || len(source.calls) != 2 {
				t.Fatalf("did not advance in same cycle: item=%v calls=%v err=%v", item, source.calls, err)
			}
			names, err := store.Load(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if allUnsupported && len(names) != 0 {
				t.Fatalf("unsupported item entered history: %v", names)
			}
			if !allUnsupported && (len(names) != 1 || names[0] != "2.gif") {
				t.Fatalf("wrong history: %v", names)
			}
		})
	}
}

func TestRotationDoesNotSkipNetworkFailures(t *testing.T) {
	store, err := history.New(filepath.Join(t.TempDir(), "history.json"))
	if err != nil {
		t.Fatal(err)
	}
	pages := &retryPages{items: []gallery.Item{{ID: 1}, {ID: 2}}}
	failure := errors.New("network timeout")
	source := &retrySource{pages: pages, failure: failure}
	runner := &Runner{selection: NewSelector(store, pages, nil), source: source}
	if _, _, err := runner.nextAnimation(context.Background()); !errors.Is(err, failure) {
		t.Fatalf("network error lost: %v", err)
	}
	if len(source.calls) != 1 || len(pages.items) != 2 {
		t.Fatal("network failure was treated as unsupported")
	}
	names, err := store.Load(context.Background())
	if err != nil || len(names) != 0 {
		t.Fatalf("failed download remained in history: %v, %v", names, err)
	}
}
