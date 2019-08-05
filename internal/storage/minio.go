package storage

import (
	"errors"
	"io"

	minio "github.com/minio/minio-go"
)

// MinIO is minio storage
type MinIO struct {
	Client     *minio.Client
	BucketName string
}

// RetreiveFile implements S
func (m MinIO) RetreiveFile(ident string) (io.ReadCloser, error) {
	if exist, _ := m.ExistFile(ident); !exist {
		return nil, errors.New("file does not exist")
	}
	return m.Client.GetObject(m.BucketName, ident, minio.GetObjectOptions{})
}

// WriteFile implements S
func (m MinIO) WriteFile(ident string, reader io.Reader) error {
	_, err := m.Client.PutObject(m.BucketName, ident, reader, -1, minio.PutObjectOptions{})
	return err
}

// ExistFile implements S
func (m MinIO) ExistFile(ident string) (bool, error) {
	_, err := m.Client.StatObject(m.BucketName, ident, minio.StatObjectOptions{})
	return err == nil, nil
}

// RenameFile implements S
func (m MinIO) RenameFile(origIdent string, nextIdent string) error {
	src := minio.NewSourceInfo(m.BucketName, origIdent, nil)
	dst, err := minio.NewDestinationInfo(m.BucketName, nextIdent, nil, nil)
	if err != nil {
		return err
	}
	if err := m.Client.CopyObject(dst, src); err != nil {
		return err
	}
	return m.Client.RemoveObject(m.BucketName, origIdent)
}

// StatFile implements S
func (m MinIO) StatFile(ident string) (*FStat, error) {
	stat, err := m.Client.StatObject(m.BucketName, ident, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &FStat{
		Length: stat.Size,
	}, nil
}
