package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

type BinaryHandler struct {
	Storage storage.S
}

// Type implements H
func (h *BinaryHandler) Type() string {
	return "binary"
}

// SizeLimit implements H
func (h *BinaryHandler) SizeLimit() int64 {
	return 50 << 20
}

// WriteData implements H
func (h *BinaryHandler) WriteData(r io.Reader, mime string, param map[string]string) (ident string, err error) {
	identEntropy := make([]byte, 512/8)
	_, err = rand.Read(identEntropy)
	if err != nil {
		return "", errors.New("crypto/rand failed to read")
	}
	identBytes := bytes.NewBuffer([]byte{})
	enc := base64.NewEncoder(base64.StdEncoding, identBytes)
	if _, err := io.Copy(enc, bytes.NewReader(identEntropy)); err != nil {
		return "", err
	}
	ident = identBytes.String()
	err = h.Storage.WriteFile(ident, r)
	return ident, err
}

// RetrieveData implements H
func (h *BinaryHandler) RetrieveData(ident string, param map[string]string) (io.ReadCloser, string, error) {
	file, err := h.Storage.RetreiveFile(ident)
	if err != nil {
		return nil, "", err
	}
	return file, "application/octet-stream", nil
}
