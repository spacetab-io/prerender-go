package service

import (
	"context"
	"fmt"
	"time"

	"github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

type Repository interface {
	SaveData(ctx context.Context, pd *models.PageData) error
}

type Service interface {
	GetLinksForRender() ([]string, error)
	GetUrlsFromSitemaps() ([]string, error)
	GetUrlsFromLinksList() ([]string, error)
	PreparePages(links []string) ([]*models.PageData, error)

	GetPageBody(ctx context.Context, p *models.PageData) error
	RenderPages(ctx context.Context, pages []*models.PageData, maxWorkers int) error
	RenderPage(ctx context.Context, page *models.PageData, num int, total int) error

	renderBodyWithElementTrigger(ctx context.Context, p *models.PageData) (string, error)
	renderBodyWithTimeTrigger(ctx context.Context, p *models.PageData) (string, error)
	renderBodyWithConsoleTrigger(ctx context.Context, p *models.PageData) (string, error)

	PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int) string
}

type service struct {
	lastRenderedAt  *time.Time
	r               Repository
	prerenderConfig configuration.PrerenderConfig
	storageConfig   configuration.StorageConfig
}

//nolint:revive // we need it here
func NewService(r Repository, prerenderConfig configuration.PrerenderConfig, storageConfig configuration.StorageConfig) *service {
	lr := time.Now().Add(-prerenderConfig.RenderPeriod)

	return &service{
		lastRenderedAt:  &lr,
		r:               r,
		prerenderConfig: prerenderConfig,
		storageConfig:   storageConfig,
	}
}

func (s *service) PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int) string {
	var (
		renderSucceed  int
		renderFailed   int
		storingSucceed int
		storingFailed  int
		result         = `
       +-------------------+
       |     status        | 
+------+---------+---------+-------+----------
|  nn  | render  | store   | tries | page path
+------+---------+---------+-------+----------
`
	)

	const (
		statusSuccess = "success"
		statusError   = "error  "
	)

	for i, page := range pages {
		var renderStatus, storingStatus string

		if page.SuccessRender {
			renderSucceed++

			renderStatus = statusSuccess
		} else {
			renderFailed++

			renderStatus = statusError
		}

		if page.SuccessStoring {
			storingSucceed++

			storingStatus = statusSuccess
		} else {
			storingFailed++

			storingStatus = statusError
		}

		result += fmt.Sprintf("| %04d | %s | %s | %d     | %s\n", i, renderStatus, storingStatus, page.Attempts, page.FileName)
	}

	format := `TOTAL info:
 - static URI: %s
 - links: %d
 - render success: %d
 - render failed: %d
 - storing success: %d
 - storing failed: %d
 - concurrent: %d
 - duration: %s
`
	result += fmt.Sprintf(format,
		getStaticURI(s.storageConfig), len(pages), renderSucceed, renderFailed, storingSucceed, storingFailed, procs, d.String())

	return result
}

func getStaticURI(config configuration.StorageConfig) string {
	switch config.Type {
	case "local":
		return config.Local.StoragePath
	case "s3":
		return config.S3.Bucket.CDNUrl + config.S3.Bucket.Folder
	}

	return "storage uri cannot be defined"
}
