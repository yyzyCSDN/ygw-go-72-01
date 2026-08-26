package version

import "errors"

var (
	ErrNoActiveVersion = errors.New("no active version")
	ErrSuperseded      = errors.New("version superseded")
	ErrVersionExists   = errors.New("version already registered")
	ErrWriteback       = errors.New("version writeback failed")
)
