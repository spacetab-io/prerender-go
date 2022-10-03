package repository

import (
	"errors"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/repository/files"
	"github.com/spacetab-io/prerender-go/pkg/repository/s3"
	"github.com/spacetab-io/prerender-go/pkg/service"
)

const (
	LocalStorage = "local"
	S3Storage    = "s3"
)

var ErrUnknownType = errors.New("storage type is unknown or  not set")

func NewRepository(storageCfg cfg.StorageConfig) (service.Repository, error) {
	switch storageCfg.Type {
	case LocalStorage:
		return files.NewStorage(storageCfg.Local.StoragePath), nil
	case S3Storage:
		return s3.NewStorage(storageCfg.S3)
	}

	return nil, ErrUnknownType
}
