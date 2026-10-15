package handlers

import (
	"errors"

	"divoom-gallery-client/internal/playback"
)

func commandMessage(err error, fallback string) string {
	if errors.Is(err, playback.ErrBusy) {
		return "The device is busy with a manual command. Wait for it to finish."
	}
	var failure *playback.Error
	if errors.As(err, &failure) {
		switch failure.Stage {
		case "settings":
			return "Could not save autoplay settings. Settings were not changed."
		case "wait":
			return "Waiting for the device was interrupted."
		case "load":
			return "The animation is unavailable or its format is unsupported. Refresh the list."
		}
	}
	return fallback
}
