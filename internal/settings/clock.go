package settings

import (
	"errors"
	"regexp"
)

type Clock struct {
	Style int    `json:"style"`
	Color string `json:"color"`
}

var clockColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func (c Clock) Validate() error {
	if c.Style < 0 || c.Style > 15 || !clockColor.MatchString(c.Color) {
		return errors.New("Choose a style from 0 to 15 and an RGB color in #RRGGBB format.")
	}
	return nil
}
