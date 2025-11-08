package file

import (
	"context"

	"github.com/IlyaChgn/voblako/internal/models"
)

type MetadataStorage interface {
	GetFilesList(ctx context.Context, ownerID uint, options models.FilesListOptions) ([]*models.FileMetadata, error)
	GetMetadata(ctx context.Context, id string) (*models.FileMetadata, error)

	UploadMetadata(ctx context.Context, ownerID uint, filename, contentType string,
		size int64) (*models.FileMetadata, error)
	UpdateFilename(ctx context.Context, id string, filename string) error
	UpdateSize(ctx context.Context, id string, size int64) error
	DeleteFile(ctx context.Context, id string) error
}

type ObjectStorage interface {
	UploadFile(ctx context.Context, key, contentType string, file []byte, size int64) error

	GetFile(ctx context.Context, key string) ([]byte, error)
}

type FileUsecases interface {
	UploadFile(ctx context.Context, ownerID uint, data *models.GeneralFileData) (*models.FileMetadata, error)
	GetFilesList(ctx context.Context, ownerID uint, options models.FilesListOptions) ([]*models.FileMetadata, error)
	GetFile(ctx context.Context, userID uint, id string) (*models.GeneralFileData, error)
	GetMetadata(ctx context.Context, userID uint, id string) (*models.FileMetadata, error)
	UpdateFile(ctx context.Context, userID uint, id string, file []byte, size int64) error
	UpdateFilename(ctx context.Context, userID uint, id string, filename string) error
	DeleteFile(ctx context.Context, userID uint, id string) error
}
