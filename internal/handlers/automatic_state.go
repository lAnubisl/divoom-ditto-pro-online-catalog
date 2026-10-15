package handlers

import (
	"strconv"
	"time"
)

type automaticState struct {
	Enabled, Running   bool
	Category           int
	Album              int64
	Error, CSRF        string
	Name               string
	Interval, Duration time.Duration
}

func (s automaticState) SourceKind() string {
	if s.Album > 0 {
		return "Album"
	}
	return "Category"
}

func durationText(d time.Duration) string {
	if d%time.Minute == 0 {
		return strconv.FormatInt(int64(d/time.Minute), 10) + " min"
	}
	return strconv.FormatFloat(d.Seconds(), 'f', -1, 64) + " s"
}
func (s automaticState) IntervalText() string { return durationText(s.Interval) }
func (s automaticState) DurationText() string { return durationText(s.Duration) }
func (s automaticState) IntervalMinutes() string {
	return strconv.FormatFloat(s.Interval.Minutes(), 'f', -1, 64)
}
func (s automaticState) DurationSeconds() string {
	return strconv.FormatFloat(s.Duration.Seconds(), 'f', -1, 64)
}
func (s automaticState) CategoryName() string {
	if s.Name != "" {
		return s.Name
	}
	return strconv.Itoa(s.Category)
}

// automationState adapts business state to the existing template model.
func (s *Server) automationState() automaticState {
	prefs, err := s.settings.Read()
	if err != nil {
		return automaticState{Error: err.Error(), CSRF: s.csrf}
	}
	runtime := s.runner.State()
	interval, duration := prefs.Timing()
	state := automaticState{Enabled: prefs.Enabled, Running: runtime.Running, Error: runtime.Error, CSRF: s.csrf, Interval: interval, Duration: duration}
	if prefs.Source == "album" {
		state.Album = prefs.SourceID
		state.Name = prefs.SourceName
	} else {
		state.Category = int(prefs.SourceID)
		for _, c := range s.categories {
			if c.ID == state.Category {
				state.Name = c.Name
			}
		}
	}
	return state
}
