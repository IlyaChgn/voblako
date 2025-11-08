package repository

import (
	"context"
	"fmt"

	"github.com/IlyaChgn/voblako/internal/models"
	"github.com/google/uuid"
)

func (s *metadataStorage) UploadMetadata(
	ctx context.Context, ownerID uint, filename, contentType string, size int64,
) (*models.FileMetadata, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var meta models.FileMetadata
	id := uuid.NewString()
	key := fmt.Sprintf("%d/%s", ownerID, id)

	line := tx.QueryRow(ctx, UploadMetadataQuery, id, ownerID, filename, contentType, size, key)
	if err := line.Scan(&meta.UUID, &meta.OwnerID, &meta.Filename, &meta.ContentType, &meta.Size,
		&meta.UploadTime, &meta.UpdateTime, &meta.StorageKey, &meta.IsDeleted, &meta.DeletedTime); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &meta, nil
}

func (s *metadataStorage) UpdateFilename(ctx context.Context, id string, filename string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, UpdateFilenameQuery, id, filename)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *metadataStorage) UpdateSize(ctx context.Context, id string, size int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, UpdateSizeQuery, id, size)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *metadataStorage) DeleteFile(ctx context.Context, id string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, DeleteFileQuery, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
