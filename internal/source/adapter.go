// Package source supplies catalog pages and cached animations.
package source

import (
	"context"
	"divoom-gallery-client/internal/animation"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/media"
	"divoom-gallery-client/internal/settings"
	"errors"
	"fmt"
)

type API interface {
	FetchCategoryPage(context.Context, int, int) (gallery.Page, error)
	FetchAlbumPage(context.Context, int64, int) (gallery.Page, error)
	FetchItem(context.Context, gallery.Item) ([]byte, error)
}

var ErrEmptyGallery = errors.New("gallery has no selectable items")

type Adapter struct {
	settings   *settings.Store
	candidates chan struct{}
	api        API
	media      *media.Store
	gate       chan struct{}
	lastPage   *cachedPage
}

type cachedPage struct {
	kind     string
	sourceID int64
	page     int
	result   gallery.Page
}

func New(api API, media *media.Store, prefs *settings.Store) *Adapter {
	return &Adapter{settings: prefs, candidates: make(chan struct{}, 1), api: api, media: media, gate: make(chan struct{}, 1)}
}

func (a *Adapter) acquire(ctx context.Context) error {
	select {
	case a.gate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// NextAnimation advances the persisted cursor and loads that item's animation.
// A selection or loading failure ends the request; the next call continues
// after the attempted item instead of retrying it immediately.
func (a *Adapter) NextAnimation(ctx context.Context) ([]byte, error) {
	item, err := a.nextCandidate(ctx)
	if err != nil {
		return nil, fmt.Errorf("select animation: %w", err)
	}
	data, err := a.Animation(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("load animation %d: %w", item.ID, err)
	}
	return data, nil
}

// nextCandidate advances the persisted source cursor before returning an item.
// Reaching the remote end resets the cursor and continues from page 1 in the
// same call. ErrEmptyGallery means a complete scan from page 1 found no
// selectable items, rather than simply reaching the end of a nonempty source.
func (a *Adapter) nextCandidate(ctx context.Context) (gallery.Item, error) {
	select {
	case a.candidates <- struct{}{}:
		defer func() { <-a.candidates }()
	case <-ctx.Done():
		return gallery.Item{}, ctx.Err()
	}
	prefs, err := a.settings.Read()
	if err != nil {
		return gallery.Item{}, err
	}
	page, previous := max(1, prefs.SourcePage), prefs.PreviousItemID
	wrapped := false
	for page <= 10000 {
		if err := ctx.Err(); err != nil {
			return gallery.Item{}, err
		}
		result, err := a.fetchPage(ctx, prefs.Source, prefs.SourceID, page)
		if err != nil {
			return gallery.Item{}, fmt.Errorf("gallery page %d: %w", page, err)
		}
		// Locate the previous item before selecting the next eligible metadata.
		start := 0
		if previous != 0 {
			for i, item := range result.Items {
				if item.ID == previous {
					start = i + 1
					break
				}
			}
		}
		for _, item := range result.Items[start:] {
			if item.Selectable(prefs.Source == "album") {
				if err := a.saveCursor(ctx, prefs, page, item.ID); err != nil {
					return gallery.Item{}, err
				}
				return item, nil
			}
		}
		if !result.HasMore() {
			if err := a.saveCursor(ctx, prefs, 0, 0); err != nil {
				return gallery.Item{}, err
			}
			if wrapped {
				return gallery.Item{}, ErrEmptyGallery
			}
			page, previous, wrapped = 1, 0, true
			continue
		}
		page++
		previous = 0
	}
	return gallery.Item{}, errors.New("gallery exceeds page limit")
}

func (a *Adapter) saveCursor(ctx context.Context, prefs settings.Values, page int, previous int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	matched := false
	err := a.settings.Update(func(current *settings.Values) {
		if current.Source == prefs.Source && current.SourceID == prefs.SourceID {
			current.SourcePage, current.PreviousItemID = page, previous
			matched = true
		}
	})
	if err != nil {
		return fmt.Errorf("save source cursor: %w", err)
	}
	if !matched {
		return errors.New("autoplay source changed during selection")
	}
	return nil
}

func (a *Adapter) fetchPage(ctx context.Context, kind string, sourceID int64, page int) (gallery.Page, error) {
	if err := ctx.Err(); err != nil {
		return gallery.Page{}, err
	}
	if cached := a.lastPage; cached != nil && cached.kind == kind && cached.sourceID == sourceID && cached.page == page {
		return cached.result, nil
	}
	result, err := a.fetchRemotePage(ctx, kind, sourceID, page)
	if err == nil {
		a.lastPage = &cachedPage{kind: kind, sourceID: sourceID, page: page, result: result}
	}
	return result, err
}

func (a *Adapter) fetchRemotePage(ctx context.Context, kind string, sourceID int64, page int) (gallery.Page, error) {
	if kind == "album" {
		return a.api.FetchAlbumPage(ctx, sourceID, page)
	}
	return a.api.FetchCategoryPage(ctx, int(sourceID), page)
}

func (a *Adapter) Animation(ctx context.Context, item gallery.Item) ([]byte, error) {
	if item.ID <= 0 {
		return nil, errors.New("invalid animation ID")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := a.acquire(ctx); err != nil {
		return nil, err
	}
	defer func() { <-a.gate }()
	if data, ok := a.media.Load(ctx, item.ID); ok {
		return data, nil
	}
	raw, err := a.api.FetchItem(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("download artwork: %w", err)
	}
	data, err := animation.Decode(raw)
	if err != nil {
		format := -1
		if len(raw) > 0 {
			format = int(raw[0])
		}
		return nil, fmt.Errorf("decode artwork (format=%d bytes=%d): %w", format, len(raw), err)
	}
	if err := a.media.Save(ctx, item.ID, data); err != nil {
		return nil, fmt.Errorf("save animation: %w", err)
	}
	return data, nil
}
