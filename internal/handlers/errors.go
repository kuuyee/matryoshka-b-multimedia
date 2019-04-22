package handlers

import "errors"

var (
	// ErrNotExist is returned when the requested file does not exist
	ErrNotExist = errors.New("file does not exist")
	// ErrMaxSizeExceeded is returned when the file exceeded its max size
	ErrMaxSizeExceeded = errors.New("max size exceeded")
)
