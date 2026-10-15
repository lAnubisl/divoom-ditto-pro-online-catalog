package gallery

import (
	"strconv"
)

type Item struct {
	ID     int64  `json:"GalleryId"`
	FileID string `json:"FileId"`
	Type   int    `json:"FileType"`
	Name   string `json:"FileName"`
	Size   int    `json:"FileSize"`
}

// Filename is the stable identity used by disk storage and playback history.
func (item Item) Filename() string { return strconv.FormatInt(item.ID, 10) + ".gif" }

// Selectable reports whether gallery metadata describes supported 16x16 work.
func (item Item) Selectable(album bool) bool {
	return item.ID > 0 && item.FileID != "" && (item.Size == 0 || item.Size == 1) &&
		(item.Type == 1 || (album && item.Type == 0))
}
