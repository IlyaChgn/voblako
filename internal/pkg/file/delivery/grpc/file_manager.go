package grpc

import (
	"context"
	"time"

	"github.com/IlyaChgn/voblako/internal/models"
	fileinterfaces "github.com/IlyaChgn/voblako/internal/pkg/file"
	"github.com/IlyaChgn/voblako/internal/pkg/file/delivery/grpc/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FileManager struct {
	protobuf.UnimplementedFileServer

	metadataStorage fileinterfaces.MetadataStorage
	objectStorage   fileinterfaces.ObjectStorage
}

func NewFileManager(
	metadataStorage fileinterfaces.MetadataStorage,
	objectStorage fileinterfaces.ObjectStorage,
) *FileManager {
	return &FileManager{
		metadataStorage: metadataStorage,
		objectStorage:   objectStorage,
	}
}

func (m *FileManager) UploadFile(ctx context.Context, r *protobuf.UploadFileRequest) (*protobuf.FileMetadata, error) {
	metadata, err := m.metadataStorage.UploadMetadata(ctx, uint(r.OwnerID), r.Filename, r.ContentType, r.Size)
	if err != nil {
		return nil, err
	}

	err = m.objectStorage.UploadFile(ctx, metadata.StorageKey, metadata.ContentType, r.Data, metadata.Size)
	if err != nil {
		return nil, err
	}

	return convertMetadata(metadata), nil
}

func (m *FileManager) GetFilesList(ctx context.Context, r *protobuf.GetFilesListRequest,
) (*protobuf.GetFilesListResponse, error) {
	resp, err := m.metadataStorage.GetFilesList(ctx, uint(r.OwnerID), models.FilesListOptions{
		Limit:       uint(r.Limit),
		Offset:      uint(r.Offset),
		WithDeleted: r.WithDeleted,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*protobuf.FileMetadata, len(resp))
	for k, v := range resp {
		list[k] = convertMetadata(v)
	}

	return &protobuf.GetFilesListResponse{Files: list}, nil
}

func (m *FileManager) GetFile(ctx context.Context, r *protobuf.GetFileRequest) (*protobuf.GetFileResponse, error) {
	meta, err := m.metadataStorage.GetMetadata(ctx, r.UUID)
	if err != nil {
		return nil, err
	} else if meta == nil || meta.OwnerID != uint(r.UserID) {
		return nil, status.Errorf(codes.PermissionDenied, models.PermissionDeniedError.Error())
	}

	file, err := m.objectStorage.GetFile(ctx, meta.StorageKey)
	if err != nil {
		return nil, err
	}

	return &protobuf.GetFileResponse{
		Filename:    meta.Filename,
		ContentType: meta.ContentType,
		Data:        file,
		Size:        meta.Size,
	}, nil
}

func (m *FileManager) GetFileMetadata(
	ctx context.Context, r *protobuf.GetFileMetadataRequest,
) (*protobuf.FileMetadata, error) {
	meta, err := m.metadataStorage.GetMetadata(ctx, r.UUID)
	if err != nil {
		return nil, err
	} else if meta == nil || meta.OwnerID != uint(r.UserID) {
		return nil, status.Errorf(codes.PermissionDenied, models.PermissionDeniedError.Error())
	}

	return convertMetadata(meta), nil
}

func (m *FileManager) UpdateFile(
	ctx context.Context, r *protobuf.UpdateFileRequest,
) (*emptypb.Empty, error) {
	meta, err := m.metadataStorage.GetMetadata(ctx, r.UUID)
	if err != nil {
		return nil, err
	} else if meta == nil || meta.OwnerID != uint(r.UserID) {
		return nil, status.Errorf(codes.PermissionDenied, models.PermissionDeniedError.Error())
	}

	err = m.objectStorage.UploadFile(ctx, meta.StorageKey, meta.ContentType, r.Data, r.Size)
	if err != nil {
		return nil, err
	}

	err = m.metadataStorage.UpdateSize(ctx, r.UUID, r.Size)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *FileManager) UpdateFilename(ctx context.Context, r *protobuf.UpdateFilenameRequest) (*emptypb.Empty, error) {
	meta, err := m.metadataStorage.GetMetadata(ctx, r.UUID)
	if err != nil {
		return nil, err
	} else if meta == nil || meta.OwnerID != uint(r.UserID) {
		return nil, status.Errorf(codes.PermissionDenied, models.PermissionDeniedError.Error())
	}

	err = m.metadataStorage.UpdateFilename(ctx, r.UUID, r.Filename)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *FileManager) DeleteFile(ctx context.Context, r *protobuf.DeleteFileRequest) (*emptypb.Empty, error) {
	meta, err := m.metadataStorage.GetMetadata(ctx, r.UUID)
	if err != nil {
		return nil, err
	} else if meta == nil || meta.OwnerID != uint(r.UserID) {
		return nil, status.Errorf(codes.PermissionDenied, models.PermissionDeniedError.Error())
	}

	err = m.metadataStorage.DeleteFile(ctx, r.UUID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func convertMetadata(m *models.FileMetadata) *protobuf.FileMetadata {
	var deletedTime time.Time
	if t := ptrTimeToProto(m.DeletedTime); t != nil {
		deletedTime = *m.DeletedTime
	}

	return &protobuf.FileMetadata{
		UUID:        m.UUID,
		OwnerID:     uint32(m.OwnerID),
		Filename:    m.Filename,
		Size:        m.Size,
		ContentType: m.ContentType,
		UploadTime:  timestamppb.New(m.UploadTime),
		UpdateTime:  timestamppb.New(m.UpdateTime),
		DeletedTime: timestamppb.New(deletedTime),
		IsDeleted:   m.IsDeleted,
	}
}

func ptrTimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
