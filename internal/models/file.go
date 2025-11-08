package models

import "time"

type FilesListOptions struct {
	Limit       uint `json:"limit"`
	Offset      uint `json:"offset"`
	WithDeleted bool `json:"with_deleted"`
}

type FileMetadata struct {
	UUID    string `json:"uuid"`
	OwnerID uint   `json:"owner_id"`

	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	StorageKey  string `json:"-"`
	IsDeleted   bool   `json:"is_deleted"`

	UploadTime  time.Time  `json:"upload_time"`
	UpdateTime  time.Time  `json:"update_time"`
	DeletedTime *time.Time `json:"deleted_time"`
}

type GeneralFileData struct {
	Filename    string
	ContentType string
	File        []byte
	Size        int64
}

type UpdateFilenameRequest struct {
	Filename string `json:"filename"`
}
