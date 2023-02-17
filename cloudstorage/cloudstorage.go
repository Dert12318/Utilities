package cloudstorage

import (
	"context"
	"io"
	"time"
)

type (
	FileOption struct {
		Object      io.Reader
		Name        string
		Size        int64
		ContentType string
	}

	UploadResponse struct {
		Bucket string
		URL    string
	}
)

type (
	CloudStorage interface {
		GetClient() interface{}
		Upload(ctx context.Context, bucketName string, makeNewBucket bool, file FileOption) (*UploadResponse, error)
		Download(ctx context.Context, bucketName, fileName string, dst io.Writer) error
		GetPreSignedURL(ctx context.Context, bucketName, fileName string, expires time.Time) (string, error)
		IsBucketExist(ctx context.Context, bucketName string) (bool, error)
		FGetObject(ctx context.Context, bucketName, objectName, filePath string) error
		DeleteObject(ctx context.Context, bucketName, objectName string) error
	}
)
