package service

import "errors"

var (
	ErrDuplicate      = errors.New("duplicate frame")
	ErrUnmappedPoint  = errors.New("unmapped point")
	ErrUnknownChannel = errors.New("unknown channel")
	ErrChannelFaulted = errors.New("channel faulted")
)
