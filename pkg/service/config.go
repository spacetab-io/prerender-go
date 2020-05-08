package service

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/spacetab-io/roastmap-go/configuration"
	"github.com/spacetab-io/roastmap-go/pkg/models"
)

func NewService(r Repository, roastmapConfig cfg.RoastmapConfig) Service {
	return &service{r, roastmapConfig}
}

type service struct {
	r   Repository
	cfg cfg.RoastmapConfig
}

func (s *service) PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int) {
	var (
		succeed int
		failed  int
	)

	for i, page := range pages {
		var status string

		if page.SuccessRender {
			succeed++

			status = "v"
		} else {
			failed++

			status = "x"
		}

		fmt.Printf("| %04d | %-100s | %s | %d |\n", i, page.FileName, status, page.Attempts)
	}

	format := `TOTAL info:
 - links: %d
 - success: %d
 - failed: %d
 - concurrent: %d
 - duration: %s
`
	fmt.Printf(format, len(pages), succeed, failed, procs, d.String())
}

type Service interface {
	GetLinksForRender() ([]string, error)
	GetUrlsFromSitemap() ([]string, error)
	GetUrlsFromLinkList() ([]string, error)
	PreparePages(links []string) ([]*models.PageData, error)

	GetPageBody(ctx context.Context, p *models.PageData) error
	RenderPages(pages []*models.PageData, maxWorkers int) error
	RenderPage(ctx context.Context, page *models.PageData, num int) error

	PrepareRenderReport(pages []*models.PageData, d time.Duration, procs int)
}

type Repository interface {
	SaveData(pd *models.PageData) error
	GzipFile() bool
}
