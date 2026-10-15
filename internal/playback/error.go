package playback

import (
	"errors"
	"fmt"
)

var ErrBusy = errors.New("a manual command is already running")

// Error identifies the failed operation so callers can present suitable feedback.
type Error struct {
	Stage string
	Err   error
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %v", e.Stage, e.Err) }

func (e *Error) Unwrap() error { return e.Err }
