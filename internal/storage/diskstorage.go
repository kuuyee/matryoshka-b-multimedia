package storage

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type Disk struct {
	BasePath string
}

func (d *Disk) joinPath(ident string) (string, error) {
	newPath := path.Join(d.BasePath, ident)
	if path.Dir(newPath) != strings.TrimSuffix(d.BasePath, string(os.PathSeparator)) {
		return "", errors.New("invalid path")
	}
	return newPath, nil
}

// RetreiveFile implements S
func (d *Disk) RetreiveFile(ident string) (io.ReadCloser, error) {
	path, err := d.joinPath(ident)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

// WriteFile implements S
func (d *Disk) WriteFile(ident string) (io.WriteCloser, error) {
	path, err := d.joinPath(ident)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}

// ExistFile implements S
func (d *Disk) ExistFile(ident string) (bool, error) {
	path, err := d.joinPath(ident)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func NewDiskStorage(basePath string) (*Disk, error) {
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, err
	}
	return &Disk{basePath}, nil
}
