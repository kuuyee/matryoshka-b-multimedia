// +build !cgo

package handlers

import (
	"errors"
	"io"

	"github.com/nfnt/resize"

	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

type ImageHandler struct {
	Storage    storage.S
	MaxSize    int64
	ResizeAlgo resize.InterpolationFunction
	KeyedMutex *KeyedRWMutex
}

// SizeLimit implements H
func (h *ImageHandler) SizeLimit() int64 {
	return h.MaxSize
}

// Type implements H
func (h *ImageHandler) Type() string {
	return "image"
}

// WriteData implements H
func (h *ImageHandler) WriteData(r io.Reader, mime string, param map[string]string) (ident string, err error) {
	return "", errors.New("cgo is required")
}

// RetrieveData implements H
func (h *ImageHandler) RetrieveData(ident string, param map[string]string) (io.ReadCloser, string, error) {
	return nil, "", errors.New("cgo is required")
}
