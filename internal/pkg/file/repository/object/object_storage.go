package repository

import (
	"bytes"
	"context"
	"io"

	fileinterfaces "github.com/IlyaChgn/voblako/internal/pkg/file"
	"github.com/minio/minio-go/v7"
)

type objectStorage struct {
	bucketName string
	client     *minio.Client
}

func NewObjectStorage(client *minio.Client, bucketName string) fileinterfaces.ObjectStorage {
	return &objectStorage{
		bucketName: bucketName,
		client:     client,
	}
}

func (s *objectStorage) UploadFile(ctx context.Context, key, contentType string, file []byte, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucketName, key, bytes.NewReader(file), size,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	return nil
}

func (s *objectStorage) GetFile(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}
