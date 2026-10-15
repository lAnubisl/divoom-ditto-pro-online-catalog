package settings

import (
	"errors"
	"time"
)

// Values preserves the existing version 2 settings file format.
type Values struct {
	Version    int    `json:"version"`
	Enabled    bool   `json:"enabled"`
	Source     string `json:"source"`
	SourceID   int64  `json:"source_id"`
	SourceName string `json:"source_name,omitempty"`
	Interval   string `json:"interval"`
	Duration   string `json:"duration"`
	StyleID    int    `json:"style_id"`
	StyleColor string `json:"style_color"`
}

func Defaults(category int, interval, duration time.Duration) Values {
	return Values{Version: 2, Enabled: true, Source: "category", SourceID: int64(category), Interval: interval.String(), Duration: duration.String(), StyleColor: "#00FF00"}
}

func (v Values) Clock() Clock { return Clock{Style: v.StyleID, Color: v.StyleColor} }

func (v Values) Timing() (time.Duration, time.Duration) {
	interval, _ := time.ParseDuration(v.Interval)
	duration, _ := time.ParseDuration(v.Duration)
	return interval, duration
}

func (v Values) validate() error {
	interval, e1 := time.ParseDuration(v.Interval)
	duration, e2 := time.ParseDuration(v.Duration)
	validSource := (v.Source == "category" && v.SourceID >= 0 && v.SourceID <= 255) || (v.Source == "album" && v.SourceID > 0)
	if err := v.Clock().Validate(); err != nil {
		return err
	}
	if v.Version != 2 || !validSource || e1 != nil || e2 != nil || interval <= 0 || duration <= 0 || duration >= interval {
		return errors.New("invalid autoplay settings: require version 2, valid source and 0 < duration < interval")
	}
	return nil
}
