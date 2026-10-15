package gallery

const PageSize = 20

type Page struct{ Items []Item }

// The API can return short pages followed by more items. Only an empty page
// marks the end of the source.
func (p Page) HasMore() bool { return len(p.Items) > 0 }
