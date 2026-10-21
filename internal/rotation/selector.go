// Package rotation selects candidates and owns automatic playback history.
package rotation

import (
	"context"
	"divoom-gallery-client/internal/gallery"
	"errors"
	"fmt"
	"log"
)

const HistoryLimit = 1000

var ErrEmptyGallery = errors.New("gallery has no selectable items")

// Selector is used sequentially by the runner.
type Selector struct {
	history  History
	source   PageSource
	metrics  Metrics
	previous []string
	pending  string
}

func NewSelector(history History, source PageSource, metrics Metrics) *Selector {
	return &Selector{history: history, source: source, metrics: metrics}
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
	r.reportHistorySize(ctx, len(next))
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
	r.reportHistorySize(ctx, len(r.previous))
	r.pending = ""
	r.previous = nil
	return nil
}

func (r *Selector) reportHistorySize(ctx context.Context, size int) {
	if r.metrics != nil {
		if err := r.metrics.WriteHistorySize(ctx, size); err != nil {
			log.Printf("report history size: %v", err)
		}
	}
}
