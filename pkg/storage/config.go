package storage

import (
	"errors"
	"strings"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/service"
	"github.com/spacetab-io/prerender-go/pkg/storage/files"
	"github.com/spacetab-io/prerender-go/pkg/storage/s3"
)

const (
	LocalStorage = "local"
	S3Storage    = "s3"
)

func NewStorage(storageCfg cfg.StorageConfig) (service.Repository, error) {
	switch storageCfg.Type {
	case LocalStorage:
		return files.NewStorage(strings.TrimRight(storageCfg.Local.StoragePath, "/")), nil
	case S3Storage:
		//return bucket.NewStorage(storageCfg.S3)
		return s3.NewStorage(storageCfg.S3), nil
	}

	return nil, errors.New("storage type is unknown or  not set")
}
