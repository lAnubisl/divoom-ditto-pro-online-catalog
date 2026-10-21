package gallery

const PageSize = 20

type Page struct {
	Items []Item
	more  bool
}

// The API can return short pages followed by more items. Only an empty page
// marks the end of the source.
// Keep the original page boundary even if filtering removes every item.
func (p Page) HasMore() bool { return p.more || len(p.Items) > 0 }
