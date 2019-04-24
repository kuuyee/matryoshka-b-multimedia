package storage

import (
	"fmt"
	"net/url"
	"strings"

	minio "github.com/minio/minio-go"

	"github.com/kuuyee/matryoshka-b-multimedia/conf"
)

// LoadStorage loads storage handler from parsed server conf
func LoadStorage() (storageH S, err error) {
	serverConf := conf.GetParsed()

	switch serverConf.Storage.Mode {
	case "disk":
		var err error
		storageH, err = NewDiskStorage(serverConf.Storage.Path)
		if err != nil {
			return nil, err
		}
		return storageH, nil
	case "minio":
		minioURL, err := url.Parse(serverConf.Storage.Path)
		if err != nil {
			return nil, err
		}
		accessKey := minioURL.User.Username()
		secretKey, _ := minioURL.User.Password()
		useSSL := minioURL.Scheme == "https"
		bucket := strings.TrimPrefix(minioURL.Path, "/")

		minioClient, err := minio.New(minioURL.Host, accessKey, secretKey, useSSL)
		if err != nil {
			return nil, err
		}

		if hasBucket, err := minioClient.BucketExists(bucket); err != nil {
			return nil, err
		} else if !hasBucket {
			err = minioClient.MakeBucket(bucket, "")
			if err != nil {
				return nil, err
			}
		}
		return &MinIO{
			Client:     minioClient,
			BucketName: bucket,
		}, nil
	default:
		return nil, fmt.Errorf("unknown storage mode: %s", serverConf.Storage.Mode)
	}
}
