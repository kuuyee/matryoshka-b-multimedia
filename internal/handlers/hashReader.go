package handlers

import (
	"encoding/hex"
	"hash"
	"io"
)

type hashReader struct {
	Hash hash.Hash
	R    io.Reader
}

// Read implements io.Reader
func (h hashReader) Read(p []byte) (n int, err error) {
	n, err = h.R.Read(p)
	if err != nil {
		return
	}
	_, err = h.Hash.Write(p[:n])
	if err != nil {
		return
	}
	return
}

func (h hashReader) Sum() []byte {
	return h.Hash.Sum(nil)
}

func (h hashReader) SumHex() string {
	return hex.EncodeToString(h.Sum())
}
