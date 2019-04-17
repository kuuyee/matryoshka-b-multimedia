package storage

import "io"

type S interface {
	RetreiveFile(ident string) (io.ReadCloser, error)
	WriteFile(ident string) (io.WriteCloser, error)
	ExistFile(ident string) (bool, error)
}
