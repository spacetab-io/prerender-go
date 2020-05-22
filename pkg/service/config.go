package service

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/models"
)

func NewService(r Repository, prerenderConfig cfg.PrerenderConfig) Service {
	return &service{r, prerenderConfig}
}

type service struct {
	r   Repository
	cfg cfg.PrerenderConfig
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

		fmt.Printf("| %04d | %s | %d | %s\n", i, status, page.Attempts, page.FileName)
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
