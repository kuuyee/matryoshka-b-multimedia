package storage

import "io"

// S is storage interface
type S interface {
	RetreiveFile(ident string) (io.ReadCloser, error)
	WriteFile(ident string, reader io.Reader) error
	ExistFile(ident string) (bool, error)
	StatFile(ident string) (*FStat, error)
	RenameFile(origIdent string, nextIdent string) error
}

// FStat is file stat result
type FStat struct {
	Length int64
}
