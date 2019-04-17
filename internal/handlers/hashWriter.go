package handlers

import (
	"encoding/hex"
	"hash"
	"io"
)

type hashWriter struct {
	Hash hash.Hash
	W    io.Writer
}

func (h hashWriter) Write(b []byte) (l int, err error) {
	l, err = h.W.Write(b)
	if err != nil {
		return
	}
	_, err = h.Hash.Write(b[:l])
	if err != nil {
		return
	}
	return
}

func (h hashWriter) Sum() []byte {
	return h.Hash.Sum(nil)
}

func (h hashWriter) SumHex() string {
	return hex.EncodeToString(h.Sum())
}
