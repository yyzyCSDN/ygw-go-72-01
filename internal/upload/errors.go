package upload

import "errors"

var (
	ErrNoSnapshot = errors.New("missing upload snapshot")
	ErrQueueFull  = errors.New("upload queue full")
	ErrWriteFail  = errors.New("upload write failed")
)
