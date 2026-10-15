package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) automationStatus(w http.ResponseWriter, r *http.Request) {
	s.render(w, "automation-state", s.automationState())
	s.render(w, "automation-oob", s.automationState())
}

func (s *Server) automationWidget(w http.ResponseWriter, r *http.Request) {
	s.render(w, "automation-widget", s.automationState())
}

func (s *Server) settingsPage(w http.ResponseWriter, r *http.Request) {
	data := pageData{Categories: s.categories, CSRF: s.csrf, Automatic: s.automationState()}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	var err error
	data.Albums, err = s.galleryAPI.ListAlbums(ctx)
	if err != nil {
		data.Error = "Could not load albums. Reload the page and try again."
	}
	data.Automatic = s.automationState()
	s.render(w, "autoplay", data)
}
func (s *Server) resumeAutoplay(w http.ResponseWriter, r *http.Request) {
	if !s.savedAutoplay(w, s.playback.Resume()) {
		return
	}
	s.automationStatus(w, r)
}
func (s *Server) pauseAutoplay(w http.ResponseWriter, r *http.Request) {
	if !s.savedAutoplay(w, s.playback.Pause()) {
		return
	}
	s.automationStatus(w, r)
}
func (s *Server) setAutoplayCategory(w http.ResponseWriter, r *http.Request) {
	category, err := strconv.Atoi(r.Form.Get("category"))
	valid := false
	for _, c := range s.categories {
		if c.ID == category {
			valid = true
		}
	}
	if err != nil || !valid {
		http.Error(w, "Invalid category", 400)
		return
	}
	if !s.savedAutoplay(w, s.playback.SetSource("category", int64(category), "")) {
		return
	}
	s.automationStatus(w, r)
}

func (s *Server) setAutoplayTiming(w http.ResponseWriter, r *http.Request) {
	minutes, intervalErr := strconv.ParseInt(r.Form.Get("interval"), 10, 64)
	seconds, durationErr := strconv.ParseInt(r.Form.Get("duration"), 10, 64)
	message := ""
	if intervalErr != nil || durationErr != nil || minutes < 1 || seconds < 1 ||
		minutes > (1<<63-1)/int64(time.Minute) || seconds > (1<<63-1)/int64(time.Second) {
		message = "Enter a positive whole number of minutes and seconds."
	} else if time.Duration(seconds)*time.Second >= time.Duration(minutes)*time.Minute {
		message = "Display duration must be shorter than the interval."
	}
	if message == "" {
		if err := s.playback.SetTiming(time.Duration(minutes)*time.Minute, time.Duration(seconds)*time.Second); err != nil {
			log.Printf("save autoplay timing: %v", err)
			message = commandMessage(err, "Could not save autoplay settings. Settings were not changed.")
		}
	}
	state := s.automationState()
	s.render(w, "autoplay-timing", struct {
		Automatic automaticState
		Error     string
		Saved     bool
	}{state, message, message == ""})
	s.render(w, "automation-state-oob", state)
	s.render(w, "automation-oob", state)
}

func (s *Server) savedAutoplay(w http.ResponseWriter, err error) bool {
	if err == nil {
		return true
	}
	log.Printf("save autoplay settings: %v", err)
	s.status(w, commandMessage(err, "Could not save autoplay settings. Settings were not changed."))
	return false
}

func (s *Server) setAutoplayAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.Form.Get("album"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album", 400)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	albums, err := s.galleryAPI.ListAlbums(ctx)
	if err != nil {
		s.status(w, "Could not validate the selected album.")
		return
	}
	for _, a := range albums {
		if a.ID == id {
			if !s.savedAutoplay(w, s.playback.SetSource("album", id, a.Name)) {
				return
			}
			s.automationStatus(w, r)
			return
		}
	}
	http.Error(w, "Unknown album", 400)
}
