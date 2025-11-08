package repository

import (
	fileinterfaces "github.com/IlyaChgn/voblako/internal/pkg/file"
	"github.com/IlyaChgn/voblako/internal/pkg/server/dbinit"
)

type metadataStorage struct {
	pool dbinit.PostgresPool
}

func NewMetadataStorage(pool dbinit.PostgresPool) fileinterfaces.MetadataStorage {
	return &metadataStorage{pool: pool}
}
