package handlers

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"

	"github.com/h2non/filetype"
	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

type AudioHandler struct {
	Storage    storage.S
	KeyedMutex *KeyedRWMutex
}

// Type implements H
func (h *AudioHandler) Type() string {
	return "audio"
}

// SizeLimit implements H
func (h *AudioHandler) SizeLimit() int64 {
	return 512 << 10
}

// WriteData implements H
func (h *AudioHandler) WriteData(r io.Reader, mime string, param map[string]string) (ident string, err error) {
	rawAudioData, err := ioutil.ReadAll(io.LimitReader(r, h.SizeLimit()))
	if err != nil {
		return "", err
	}

	if typ, err := filetype.Audio(rawAudioData); err != nil || typ.Extension != "ogg" {
		return "", errors.New("uploaded audio file is not ogg")
	}

	processedAudioData := bytes.NewBuffer([]byte{})
	hasher := hashWriter{
		Hash: sha256.New(),
		W:    processedAudioData,
	}

	if _, err := io.Copy(hasher, bytes.NewBuffer(rawAudioData)); err != nil {
		return "", err
	}

	ident = hasher.SumHex() + ".ogg"

	h.KeyedMutex.GetMutex(ident).Lock()
	defer h.KeyedMutex.GetMutex(ident).Unlock()

	if err := h.Storage.WriteFile(ident, processedAudioData); err != nil {
		return "", err
	}
	return ident, nil
}

// RetrieveData implements H
func (h *AudioHandler) RetrieveData(ident string, param map[string]string) (io.ReadCloser, string, error) {
	file, err := h.Storage.RetreiveFile(ident)
	if err != nil {
		return nil, "", err
	}

	return file, "audio/ogg", nil
}
