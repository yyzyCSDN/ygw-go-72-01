package channel

import "errors"

var (
	ErrUnknownChannel = errors.New("unknown channel")
	ErrInvalidState   = errors.New("invalid state transition")
	ErrSessionClosed  = errors.New("session closed")
	ErrQueueFull      = errors.New("channel queue full")
	ErrEmptyBatch     = errors.New("empty batch")
	ErrStaleSync      = errors.New("stale sync record")
)
