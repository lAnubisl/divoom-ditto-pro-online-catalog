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

type Adapter struct {
	settings *settings.Store
	api      API
	media    *media.Store
	gate     chan struct{}
}

func New(prefs *settings.Store, api API, media *media.Store) *Adapter {
	return &Adapter{settings: prefs, api: api, media: media, gate: make(chan struct{}, 1)}
}

func (a *Adapter) acquire(ctx context.Context) error {
	select {
	case a.gate <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Pages captures the configured source for one rotation decision. The returned
// function supplies supported items and reports the end of the remote source
// separately, so a page of unsupported items does not end traversal.
func (a *Adapter) Pages() (func(context.Context, int) ([]gallery.Item, bool, error), error) {
	prefs, err := a.settings.Read()
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, page int) ([]gallery.Item, bool, error) {
		var result gallery.Page
		var err error
		if prefs.Source == "album" {
			result, err = a.api.FetchAlbumPage(ctx, prefs.SourceID, page)
		} else {
			result, err = a.api.FetchCategoryPage(ctx, int(prefs.SourceID), page)
		}
		if err != nil {
			return nil, false, err
		}
		items := make([]gallery.Item, 0, len(result.Items))
		for _, item := range result.Items {
			if item.Selectable(prefs.Source == "album") {
				items = append(items, item)
			}
		}
		return items, !result.HasMore(), nil
	}, nil
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
