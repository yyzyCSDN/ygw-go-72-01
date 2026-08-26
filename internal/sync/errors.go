package timesync

import "errors"

var (
	ErrStaleResult = errors.New("stale timesync result")
	ErrNoSource    = errors.New("missing timesource")
)
