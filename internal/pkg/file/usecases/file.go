package usecases

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/IlyaChgn/voblako/internal/models"
	fileinterfaces "github.com/IlyaChgn/voblako/internal/pkg/file"
	"github.com/IlyaChgn/voblako/internal/pkg/file/delivery/grpc/protobuf"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fileUsecases struct {
	client protobuf.FileClient
}

func NewFileUsecases(client protobuf.FileClient) fileinterfaces.FileUsecases {
	return &fileUsecases{client: client}
}

func (uc *fileUsecases) UploadFile(ctx context.Context, ownerID uint,
	data *models.GeneralFileData) (*models.FileMetadata, error) {
	if utf8.RuneCountInString(data.Filename) < 0 || utf8.RuneCountInString(data.Filename) > 50 {
		return nil, models.InvalidInputError
	}

	metadata, err := uc.client.UploadFile(ctx, &protobuf.UploadFileRequest{
		OwnerID:     uint32(ownerID),
		Filename:    data.Filename,
		Data:        data.File,
		ContentType: data.ContentType,
		Size:        data.Size,
	})
	if err != nil {
		return nil, err
	}

	return convertMetadata(metadata), nil
}

func (uc *fileUsecases) GetFilesList(ctx context.Context, ownerID uint,
	options models.FilesListOptions) ([]*models.FileMetadata, error) {
	if int(options.Limit) < 0 || int(options.Offset) < 0 {
		return nil, models.InvalidInputError
	}

	resp, err := uc.client.GetFilesList(ctx, &protobuf.GetFilesListRequest{
		OwnerID:     uint32(ownerID),
		Limit:       uint32(options.Limit),
		Offset:      uint32(options.Offset),
		WithDeleted: options.WithDeleted,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*models.FileMetadata, len(resp.Files))
	for k, v := range resp.Files {
		list[k] = convertMetadata(v)
	}

	return list, nil
}

func (uc *fileUsecases) GetFile(ctx context.Context, userID uint, id string) (*models.GeneralFileData, error) {
	err := uuid.Validate(id)
	if err != nil {
		return nil, models.InvalidInputError
	}

	fileData, err := uc.client.GetFile(ctx, &protobuf.GetFileRequest{
		UUID:   id,
		UserID: uint32(userID),
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			return nil, models.PermissionDeniedError
		}

		return nil, err
	}

	return &models.GeneralFileData{
		Filename:    fileData.Filename,
		ContentType: fileData.ContentType,
		File:        fileData.Data,
		Size:        fileData.Size,
	}, nil
}

func (uc *fileUsecases) GetMetadata(ctx context.Context, userID uint, id string) (*models.FileMetadata, error) {
	err := uuid.Validate(id)
	if err != nil {
		return nil, models.InvalidInputError
	}

	metadata, err := uc.client.GetFileMetadata(ctx, &protobuf.GetFileMetadataRequest{
		UUID:   id,
		UserID: uint32(userID),
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			return nil, models.PermissionDeniedError
		}

		return nil, err
	}

	return convertMetadata(metadata), nil
}

func (uc *fileUsecases) UpdateFile(ctx context.Context, userID uint, id string, file []byte, size int64) error {
	err := uuid.Validate(id)
	if err != nil {
		return models.InvalidInputError
	}

	_, err = uc.client.UpdateFile(ctx, &protobuf.UpdateFileRequest{
		UUID:   id,
		UserID: uint32(userID),
		Data:   file,
		Size:   size,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			return models.PermissionDeniedError
		}

		return err
	}

	return nil
}

func (uc *fileUsecases) UpdateFilename(ctx context.Context, userID uint, id string, filename string) error {
	if utf8.RuneCountInString(filename) < 0 || utf8.RuneCountInString(filename) > 50 {
		return models.InvalidFilenameError
	}

	err := uuid.Validate(id)
	if err != nil {
		return models.InvalidInputError
	}

	_, err = uc.client.UpdateFilename(ctx, &protobuf.UpdateFilenameRequest{
		UUID:     id,
		UserID:   uint32(userID),
		Filename: filename,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			return models.PermissionDeniedError
		}

		return err
	}

	return nil
}

func (uc *fileUsecases) DeleteFile(ctx context.Context, userID uint, id string) error {
	err := uuid.Validate(id)
	if err != nil {
		return models.InvalidInputError
	}

	_, err = uc.client.DeleteFile(ctx, &protobuf.DeleteFileRequest{
		UUID:   id,
		UserID: uint32(userID),
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			return models.PermissionDeniedError
		}

		return err
	}

	return nil
}

func convertMetadata(meta *protobuf.FileMetadata) *models.FileMetadata {
	return &models.FileMetadata{
		UUID:        meta.UUID,
		OwnerID:     uint(meta.OwnerID),
		Filename:    meta.Filename,
		ContentType: meta.ContentType,
		Size:        meta.Size,
		IsDeleted:   meta.IsDeleted,
		UploadTime:  meta.UploadTime.AsTime(),
		UpdateTime:  meta.UpdateTime.AsTime(),
		DeletedTime: protoToPtrTime(meta.DeletedTime),
	}
}

func protoToPtrTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	tm := t.AsTime()
	return &tm
}
