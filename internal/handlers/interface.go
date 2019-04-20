package handlers

import (
	"io"
)

type H interface {
	// return handler type
	Type() string

	// max size in bytes
	SizeLimit() int64

	WriteData(r io.Reader, mime string, param map[string]string) (ident string, err error)
	RetrieveData(ident string, param map[string]string) (io.ReadCloser, string, error)
}
