package parse

import "errors"

var (
	ErrEmptyFrame          = errors.New("empty frame")
	ErrBadFrame            = errors.New("malformed frame")
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
	ErrFrameTooLong        = errors.New("frame too long")
)
