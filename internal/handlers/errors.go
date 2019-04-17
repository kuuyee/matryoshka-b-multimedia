package handlers

import "errors"

var (
	ErrNotExist        = errors.New("file does not exist")
	ErrMaxSizeExceeded = errors.New("max size exceeded")
)
