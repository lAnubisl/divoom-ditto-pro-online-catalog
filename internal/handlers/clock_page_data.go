package handlers

import (
	"divoom-gallery-client/internal/settings"
)

type clockPageData struct {
	Automatic        automaticState
	CSRF, Error      string
	Settings         settings.Clock
	Styles           []int
	Available, Saved bool
}
