package aliyun_oss

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/Dert12318/Utilities/cloudstorage"
)

type (
	Option struct {
		RegionName      string
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
	}

	cloudStorage struct {
		client     *oss.Client
		regionName string
	}
)

func NewCloudStorage(opt Option) (cloudstorage.CloudStorage, error) {
	// Initialize Aliyun OSS client object.
	ossClient, err := oss.New(opt.Endpoint, opt.AccessKeyID, opt.SecretAccessKey)

	if err != nil {
		log.Fatalln(err)
	}

	return &cloudStorage{
		client:     ossClient,
		regionName: opt.RegionName,
	}, nil
}

func (c *cloudStorage) GetClient() interface{} {
	return c.client
}

func (c *cloudStorage) Upload(ctx context.Context, bucketName string, makeNewBucket bool, file cloudstorage.FileOption) (*cloudstorage.UploadResponse, error) {
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	err = bucket.PutObject(file.Name, file.Object)
	if err != nil {
		return nil, err
	}

	return &cloudstorage.UploadResponse{
		Bucket: bucket.BucketName,
		URL:    file.Name,
	}, nil
}

func (c *cloudStorage) Download(ctx context.Context, bucketName, fileName string, dst io.Writer) error {
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Case 1: Download the object into ReadCloser(). The body needs to be closed
	body, err := bucket.GetObject(fileName)
	if err != nil {
		return err
	}
	defer func() { body.Close() }()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	_, err = dst.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *cloudStorage) GetPreSignedURL(ctx context.Context, bucketName, fileName string, expires time.Time) (string, error) {
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	// generate presign URL
	presignedURL, err := bucket.SignURL(fileName, oss.HTTPGet, expires.Unix())
	if err != nil {
		return "", err
	}

	return presignedURL, nil
}

func (c *cloudStorage) FGetObject(ctx context.Context, bucketName, objectName, filePath string) error {
	return errors.New("NOT IMPLEMENTED")
}

func (c *cloudStorage) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {
	//check if bucket exist
	return c.client.IsBucketExist(bucketName)
}

func (c *cloudStorage) createBucket(ctx context.Context, bucketName string) error {
	err := c.client.CreateBucket(bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (c *cloudStorage) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	bucket, err := c.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}
