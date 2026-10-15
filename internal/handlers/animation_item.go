package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"divoom-gallery-client/internal/gallery"
)

// itemToken carries the listed item across independent HTTP requests without a
// server-side metadata map. Its signature prevents changing the cache identity
// or download path independently of the item returned by the gallery.
func (s *Server) itemToken(item gallery.Item) (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(data)
	mac := hmac.New(sha256.New, []byte(s.itemKey))
	mac.Write([]byte(payload))
	return payload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (s *Server) animationItem(r *http.Request) (gallery.Item, error) {
	var item gallery.Item
	id, err := animationID(r)
	if err != nil {
		return item, err
	}
	token := r.FormValue("item")
	if len(token) > 8192 {
		return item, errors.New("invalid animation item")
	}
	payload, signature, ok := strings.Cut(token, ".")
	if !ok {
		return item, errors.New("invalid animation item")
	}
	mac := hmac.New(sha256.New, []byte(s.itemKey))
	mac.Write([]byte(payload))
	provided, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil || !hmac.Equal(provided, mac.Sum(nil)) {
		return item, errors.New("invalid animation item signature")
	}
	data, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil || json.Unmarshal(data, &item) != nil || item.ID != id {
		return gallery.Item{}, errors.New("invalid animation item")
	}
	return item, nil
}
