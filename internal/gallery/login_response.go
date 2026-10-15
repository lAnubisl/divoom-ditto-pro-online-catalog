package gallery

import (
	"encoding/json"
)

type loginResponse struct {
	Token  json.RawMessage
	UserID json.RawMessage `json:"UserId"`
}
