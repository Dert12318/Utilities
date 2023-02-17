package minio

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Dert12318/Utilities/cloudstorage"
)

type (
	Option struct {
		RegionName          string
		Endpoint            string
		AccessKeyID         string
		SecretAccessKey     string
		UseSSL              bool
		EnableObjectLocking bool
		AllowCreateBucket   bool
	}

	cloudStorage struct {
		client            *minio.Client
		regionName        string
		objectLocking     bool
		allowCreateBucket bool
	}
)

func NewCloudStorage(opt Option) (cloudstorage.CloudStorage, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(opt.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opt.AccessKeyID, opt.SecretAccessKey, ""),
		Secure: opt.UseSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &cloudStorage{
		client:        minioClient,
		regionName:    opt.RegionName,
		objectLocking: opt.EnableObjectLocking,
	}, nil
}

func (c *cloudStorage) GetClient() interface{} {
	return c.client
}

func (c *cloudStorage) Upload(ctx context.Context, bucketName string, makeNewBucket bool, file cloudstorage.FileOption) (*cloudstorage.UploadResponse, error) {
	//create bucket
	if makeNewBucket {
		if err := c.createBucket(ctx, bucketName); err != nil {
			return nil, err
		}
	}

	result, err := c.client.PutObject(ctx, bucketName, file.Name, file.Object, file.Size,
		minio.PutObjectOptions{ContentType: file.ContentType})
	if err != nil {
		return nil, err
	}
	return &cloudstorage.UploadResponse{
		Bucket: result.Bucket,
		URL:    result.Key,
	}, nil
}

func (c *cloudStorage) Download(ctx context.Context, bucketName, fileName string, dst io.Writer) error {
	object, err := c.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	defer func() { object.Close() }()

	_, err = io.Copy(dst, object)
	if err != nil {
		return err
	}

	return nil
}

func (c *cloudStorage) FGetObject(ctx context.Context, bucketName, objectName, filePath string) error {
	err := c.client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *cloudStorage) GetPreSignedURL(ctx context.Context, bucketName, fileName string, expires time.Time) (string, error) {
	preSignedURL, err := c.client.PresignedGetObject(ctx, bucketName, fileName, expires.Sub(time.Now().Add(0*time.Second)), nil)
	if err != nil {
		return "", err
	}

	return preSignedURL.String(), nil
}

func (c *cloudStorage) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {
	//check if bucket exist
	bucketExist, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return bucketExist, nil
}

func (c *cloudStorage) createBucket(ctx context.Context, bucketName string) error {
	err := c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{ObjectLocking: c.objectLocking, Region: c.regionName})
	if err != nil {
		return err
	}

	return nil
}

func (c *cloudStorage) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	opt := minio.RemoveObjectOptions{
		GovernanceBypass: false,
	}
	err := c.client.RemoveObject(ctx, bucketName, objectName, opt)
	if err != nil {
		return err
	}

	return nil
}
