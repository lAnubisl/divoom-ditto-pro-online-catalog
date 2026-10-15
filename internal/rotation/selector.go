// Package rotation selects candidates and owns automatic playback history.
package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
	"errors"
	"fmt"
)

const HistoryLimit = 1000

var ErrEmptyGallery = errors.New("gallery has no selectable items")

// Selector is used sequentially by the runner.
type Selector struct {
	history  History
	source   PageSource
	previous []string
	pending  string
}

func NewSelector(history History, source PageSource) *Selector {
	return &Selector{history: history, source: source}
}

// Next records the selected candidate. Reject removes it if retrieval or sending
// fails. Each decision starts at page one with the latest source settings.
func (r *Selector) Next(ctx context.Context) (gallery.Item, error) {
	history, err := r.history.Load(ctx)
	if err != nil {
		return gallery.Item{}, fmt.Errorf("load playback history: %w", err)
	}
	pages, err := r.source.Pages()
	if err != nil {
		return gallery.Item{}, err
	}
	known := make(map[string]bool, len(history))
	for _, name := range history {
		known[name] = true
	}
	var first gallery.Item
	for page := 1; ; page++ {
		if err := ctx.Err(); err != nil {
			return gallery.Item{}, err
		}
		items, end, err := pages(ctx, page)
		if err != nil {
			return gallery.Item{}, fmt.Errorf("gallery page %d: %w", page, err)
		}
		for _, item := range items {
			if first.ID == 0 {
				first = item
			}
			if !known[item.Filename()] {
				return r.record(ctx, history, item)
			}
		}
		if end {
			break
		}
	}
	if first.ID == 0 {
		return gallery.Item{}, ErrEmptyGallery
	}
	// Replacing history restarts the source. Retaining the whole first item avoids
	// a second page scan or any lookup of its download metadata.
	return r.record(ctx, nil, first)
}

func (r *Selector) record(ctx context.Context, history []string, item gallery.Item) (gallery.Item, error) {
	next := append(append([]string{}, history...), item.Filename())
	if len(next) > HistoryLimit {
		next = next[len(next)-HistoryLimit:]
	}
	if err := r.history.Save(ctx, next); err != nil {
		return gallery.Item{}, err
	}
	r.previous = history
	r.pending = item.Filename()
	return item, nil
}

func (r *Selector) Reject(ctx context.Context, item gallery.Item) error {
	if r.pending != item.Filename() {
		return errors.New("candidate is not pending")
	}
	if err := r.history.Save(ctx, r.previous); err != nil {
		return err
	}
	r.pending = ""
	r.previous = nil
	return nil
}
