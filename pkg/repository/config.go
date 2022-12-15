package repository

import (
	"fmt"

	"github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/errors"
	"github.com/spacetab-io/prerender-go/pkg/repository/files"
	"github.com/spacetab-io/prerender-go/pkg/repository/s3"
	"github.com/spacetab-io/prerender-go/pkg/service"
)

const (
	LocalStorage = "local"
	S3Storage    = "s3"
)

//nolint:ireturn,nolintlint // we need it here
func NewRepository(storageCfg configuration.StorageConfig) (service.Repository, error) {
	var (
		rep service.Repository
		err error
	)

	switch storageCfg.Type {
	case LocalStorage:
		rep, err = files.NewStorage(storageCfg.Local.StoragePath), nil
	case S3Storage:
		rep, err = s3.NewStorage(storageCfg.S3)
	default:
		return nil, errors.ErrUnknownType
	}

	if err != nil {
		return nil, fmt.Errorf("new repository error: %w", err)
	}

	return rep, nil
}
