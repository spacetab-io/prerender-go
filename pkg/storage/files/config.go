package files

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spacetab-io/prerender-go/pkg/models"
)

type storage struct {
	path string
}

func NewStorage(folderPath string) storage { //nolint:golint
	return storage{folderPath}
}

func (s storage) GzipFile() bool {
	return false
}

func (s storage) SaveData(pd *models.PageData) error {
	if pd == nil {
		return errors.New("nil page data")
	}

	fullPath := s.path + "/" + pd.FileName

	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("making dir error: :%v", err)
	}

	err := ioutil.WriteFile(fullPath, pd.Body, 0644)
	if err != nil {
		return fmt.Errorf("writing file error: %v", err)
	}

	// clear body to release memory
	pd.Body = nil

	return nil
}
