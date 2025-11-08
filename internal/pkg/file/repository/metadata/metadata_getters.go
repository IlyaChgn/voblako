package repository

import (
	"context"
	"errors"

	"github.com/IlyaChgn/voblako/internal/models"
	"github.com/jackc/pgx/v5"
)

func (s *metadataStorage) GetFilesList(
	ctx context.Context, ownerID uint, options models.FilesListOptions,
) ([]*models.FileMetadata, error) {
	rows, err := s.pool.Query(ctx, GetFilesListQuery, ownerID, options.WithDeleted, options.Limit, options.Offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	var list []*models.FileMetadata
	for rows.Next() {
		var meta models.FileMetadata
		if err := rows.Scan(&meta.UUID, &meta.OwnerID, &meta.Filename, &meta.ContentType, &meta.Size,
			&meta.UploadTime, &meta.UpdateTime, &meta.StorageKey, &meta.IsDeleted, &meta.DeletedTime); err != nil {
			return nil, err
		}

		list = append(list, &meta)
	}

	return list, nil
}

func (s *metadataStorage) GetMetadata(ctx context.Context, id string) (*models.FileMetadata, error) {
	var meta models.FileMetadata

	row := s.pool.QueryRow(ctx, GetMetadataQuery, id)
	if err := row.Scan(&meta.UUID, &meta.OwnerID, &meta.Filename, &meta.ContentType, &meta.Size,
		&meta.UploadTime, &meta.UpdateTime, &meta.StorageKey, &meta.IsDeleted, &meta.DeletedTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &meta, nil
}
