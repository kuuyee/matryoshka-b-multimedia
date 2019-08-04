package storage

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

// Disk is disk storage
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
func (d *Disk) WriteFile(ident string, reader io.Reader) error {
	path, err := d.joinPath(ident)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	return nil
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

// RenameFile implements S
func (d *Disk) RenameFile(origIdent string, nextIdent string) error {
	origPath, err := d.joinPath(origIdent)
	if err != nil {
		return err
	}
	nextPath, err := d.joinPath(nextIdent)
	if err != nil {
		return err
	}
	return os.Rename(origPath, nextPath)
}

// NewDiskStorage creates a new disk storage handler
func NewDiskStorage(basePath string) (*Disk, error) {
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, err
	}
	return &Disk{basePath}, nil
}
