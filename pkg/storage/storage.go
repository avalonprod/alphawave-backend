package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	EndpointURL string
	minioClient *minio.Client
}

func NewClient(endpoint, accessKeyID, secretAccessKey string) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %v", err)
	}

	return &Client{
		EndpointURL: endpoint,
		minioClient: minioClient,
	}, nil
}

func (c *Client) GetFilePresignedURL(ctx context.Context, bucketName, fileName string, expiresTime time.Duration) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reqParams := make(url.Values)

	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

	presignedURL, err := c.minioClient.PresignedGetObject(nCtx, bucketName, fileName, expiresTime, reqParams)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", presignedURL.Host, presignedURL.Path), nil

}

func (c *Client) GetFile(ctx context.Context, bucketName, fileName string) (*[]byte, error) {
	nCtx, cancel := context.WithTimeout(ctx, 70*time.Second)
	defer cancel()

	object, err := c.minioClient.GetObject(nCtx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {

		return nil, err
	}
	defer object.Close()

	objectInfo, err := object.Stat()

	if err != nil {

		return nil, err
	}
	buffer := make([]byte, objectInfo.Size)

	_, err = object.Read(buffer)

	if err != nil {

		return nil, fmt.Errorf("failed to get object. err: %w", err)
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get object. err: %w", err)
	}

	return &buffer, nil
}

func (c *Client) UploadFile(ctx context.Context, bucketName, objectName, fileName string, fileSize int64, reader io.Reader) error {
	nCtx, cancel := context.WithTimeout(ctx, 70*time.Second)
	defer cancel()

	exists, errBucketExists := c.minioClient.BucketExists(ctx, bucketName)

	if errBucketExists != nil || !exists {
		err := c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	_, err := c.minioClient.PutObject(nCtx, bucketName, objectName, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name":      fileName,
				"x-amz-acl": "public-read",
			},
			ContentType: "application/octet-stream",
		})
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}
	return nil
}

func (c *Client) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	nCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	if err := c.minioClient.RemoveObject(nCtx, bucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}
