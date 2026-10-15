package handlers

import (
	"divoom-gallery-client/internal/settings"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) clockData() clockPageData {
	d := clockPageData{Automatic: s.automationState(), CSRF: s.csrf}
	for i := 0; i < 16; i++ {
		d.Styles = append(d.Styles, i)
	}
	prefs, err := s.settings.Read()
	if err == nil {
		d.Settings = prefs.Clock()
		d.Available = true
	} else {
		d.Error = err.Error()
	}
	return d
}

func (s *Server) clockPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "clock-settings", s.clockData())
}

func (s *Server) applyClockSettings(w http.ResponseWriter, r *http.Request) {
	d := s.clockData()
	style, err := strconv.Atoi(r.Form.Get("style"))
	clock := settings.Clock{Style: style, Color: r.Form.Get("color")}
	if !d.Available {
		d.Error = "Clock settings are unavailable."
	} else if err != nil || clock.Validate() != nil {
		d.Error = "Choose a style from 0 to 15 and an RGB color in #RRGGBB format."
	} else {
		applied, err := s.playback.ApplyClock(r.Context(), clock)
		if applied {
			d.Settings = clock
		}
		if err != nil {
			log.Printf("clock settings: %v", err)
			if applied {
				d.Error = "Clock settings were applied, but could not be saved. They may reset after restart."
			} else {
				d.Error = commandMessage(err, "Could not apply clock settings. Try again.")
			}
		} else {
			d.Saved = true
		}
	}

	d.Automatic = s.automationState()
	s.render(w, "clock-form", d)
	s.render(w, "automation-oob", d.Automatic)
}

func (s *Server) restoreClock(w http.ResponseWriter, r *http.Request) {
	if err := s.playback.RestoreClock(r.Context()); err != nil {
		log.Printf("clock: %v", err)
		s.status(w, commandMessage(err, "Could not restore the clock. Try again."))
		return
	}
	s.status(w, "The server completed the clock command.")
}
