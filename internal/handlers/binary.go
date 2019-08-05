package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"

	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

var randChars = []byte("1234567890" + "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randCharStr(lenStr int) string {
	var w bytes.Buffer
	lenChars := big.NewInt(int64(len(randChars)))
	for i := 0; i < lenStr; i++ {
		index, _ := rand.Int(rand.Reader, lenChars)
		w.WriteByte(randChars[int(index.Int64())])
	}
	return w.String()
}

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
	tmpObjName := "tmp-" + randCharStr(8)
	hashObjReader := &hashReader{
		Hash: sha256.New(),
		R:    r,
	}
	err = h.Storage.WriteFile(tmpObjName, hashObjReader)
	if err != nil {
		return "", err
	}
	ident = hashObjReader.SumHex()
	err = h.Storage.RenameFile(tmpObjName, ident)
	return ident, err
}

// RetrieveData implements H
func (h *BinaryHandler) RetrieveData(ident string, param map[string]string) (io.ReadCloser, int64, string, error) {
	stat, err := h.Storage.StatFile(ident)
	if err != nil {
		return nil, 0, "", err
	}
	file, err := h.Storage.RetreiveFile(ident)
	if err != nil {
		return nil, 0, "", err
	}
	return file, stat.Length, "application/octet-stream", nil
}
