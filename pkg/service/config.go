package service

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

func NewService(r Repository, prerenderConfig cfg.PrerenderConfig, storageConfig cfg.StorageConfig) Service {
	return &service{r, prerenderConfig, storageConfig}
}

type service struct {
	r               Repository
	prerenderConfig cfg.PrerenderConfig
	storageConfig   cfg.StorageConfig
}

func (s *service) PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int) {
	var (
		renderSucceed  int
		renderFailed   int
		storingSucceed int
		storingFailed  int
	)

	fmt.Print(`
       +-------------------+
       |     status        | 
+------+---------+---------+-------+----------
|  nn  | render  | store   | tries | page path
+------+---------+---------+-------+----------
`)
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

		fmt.Printf("| %04d | %s | %s | %d     | %s\n", i, renderStatus, storingStatus, page.Attempts, page.FileName)
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
	fmt.Printf(format,
		getStaticURI(s.storageConfig), len(pages), renderSucceed, renderFailed, storingSucceed, storingFailed, procs, d.String())
}

func getStaticURI(config cfg.StorageConfig) string {
	switch config.Type {
	case "local":
		return config.Local.StoragePath
	case "s3":
		return config.S3.Bucket.CDNUrl + config.S3.Bucket.Folder
	}

	return "storage uri cannot be defined"
}

type Service interface {
	GetLinksForRender() ([]string, error)
	GetUrlsFromSitemaps() ([]string, error)
	GetUrlsFromLinksList() ([]string, error)
	PreparePages(links []string) ([]*models.PageData, error)

	GetPageBody(ctx context.Context, p *models.PageData) error
	RenderPages(pages []*models.PageData, maxWorkers int) error
	RenderPage(ctx context.Context, page *models.PageData, num int) error

	renderBodyWithElementTrigger(ctx context.Context, p *models.PageData) (string, error)
	renderBodyWithTimeTrigger(ctx context.Context, p *models.PageData) (string, error)
	renderBodyWithConsoleTrigger(ctx context.Context, p *models.PageData) (string, error)

	PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int)
}

type Repository interface {
	SaveData(pd *models.PageData) error
}
