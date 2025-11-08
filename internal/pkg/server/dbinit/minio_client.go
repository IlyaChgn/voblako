package dbinit

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioEndpoint(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func NewMinioClient(endpoint, accessKey, secretKey, bucketName string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Created bucket: %s\n", bucketName)
	} else {
		log.Printf("Bucket %s already exists\n", bucketName)
	}

	return client, nil
}
