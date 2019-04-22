package test

import (
	"io/ioutil"

	"github.com/kuuyee/matryoshka-b-multimedia/internal/storage"
)

// NewTmpStorage creates a storage handler with os temp directory
func NewTmpStorage(ident string) storage.S {
	tmpDir, err := ioutil.TempDir("", ident)
	if err != nil {
		panic(err)
	}
	d, err := storage.NewDiskStorage(tmpDir)
	if err != nil {
		panic(err)
	}
	return d
}
