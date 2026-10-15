package gallery

import (
	"encoding/json"
)

type albumEntry struct {
	ID   json.RawMessage `json:"AlbumId"`
	Name string          `json:"AlbumName"`
}
