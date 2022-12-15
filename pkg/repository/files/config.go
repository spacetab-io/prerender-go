package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacetab-io/prerender-go/pkg/errors"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

type storage struct {
	path string
}

//nolint:revive // we need it here
func NewStorage(folderPath string) *storage {
	return &storage{path: strings.TrimRight(folderPath, "/")}
}

//nolint:gomnd // permissions numbers
func (s storage) SaveData(_ context.Context, pd *models.PageData) error {
	if pd == nil {
		return errors.ErrPageIsNil
	}

	fullPath := s.path + "/" + pd.FileName

	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("making dir error: :%w", err)
	}

	err := os.WriteFile(fullPath, pd.Body, 0o600)
	if err != nil {
		return fmt.Errorf("writing file error: %w", err)
	}

	// clear body to release memory
	pd.Body = nil

	return nil
}
