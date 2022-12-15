package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/spacetab-io/prerender-go/pkg/errors"
	"github.com/spacetab-io/prerender-go/pkg/models"
	"golang.org/x/sync/semaphore"
)

func (s service) GetPageBody(ctx context.Context, p *models.PageData) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, s.prerenderConfig.WaitTimeout)
	defer cancel()

	newTabCtx, cancelNewTabCtx := chromedp.NewContext(timeoutCtx) // create new tab
	defer cancelNewTabCtx()

	var (
		body string
		err  error
	)

	switch s.prerenderConfig.WaitFor {
	case models.WaitForConsole:
		body, err = s.renderBodyWithConsoleTrigger(newTabCtx, p)
	case models.WaitForElement:
		body, err = s.renderBodyWithElementTrigger(newTabCtx, p)
	case models.WaitForTime:
		body, err = s.renderBodyWithTimeTrigger(newTabCtx, p)
	default:
		err = errors.ErrUnknownTrigger
	}

	if err != nil {
		return fmt.Errorf("get page body error: %w", err)
	}

	p.Body = []byte(body)
	p.ContentLength = len(body)
	p.Status = http.StatusOK

	if s.prerenderConfig.Page404Text != "" && strings.Contains(body, s.prerenderConfig.Page404Text) {
		p.Status = http.StatusNotFound
	}

	return nil
}

func (s service) renderBodyWithElementTrigger(ctx context.Context, p *models.PageData) (string, error) {
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.prerenderConfig.Viewport.Width, s.prerenderConfig.Viewport.Height, 1.0, false),
		chromedp.WaitReady(s.prerenderConfig.Element.GetWaitElement()),
		chromedp.WaitReady(s.prerenderConfig.Element.GetWaitElementAttr("ready")),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		return "", fmt.Errorf("renderBodyWithElementTrigger error: %w", err)
	}

	return body, nil
}

func (s service) renderBodyWithTimeTrigger(ctx context.Context, p *models.PageData) (string, error) {
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.prerenderConfig.Viewport.Width, s.prerenderConfig.Viewport.Height, 1.0, false),
		chromedp.Sleep(s.prerenderConfig.SleepTime),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		return "", fmt.Errorf("renderBodyWithTimeTrigger error: %w", err)
	}

	return body, nil
}

func (s service) renderBodyWithConsoleTrigger(ctx context.Context, p *models.PageData) (string, error) {
	gotResult := make(chan bool, 1)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			if ev.Type == runtime.APITypeLog {
				for _, arg := range ev.Args {
					if string(arg.Value) == fmt.Sprintf(`"%s"`, s.prerenderConfig.ConsoleString) {
						gotResult <- true
					}
				}
			}
		case *runtime.EventExceptionThrown:
		}
	})

	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.prerenderConfig.Viewport.Width, s.prerenderConfig.Viewport.Height, 1.0, false),
		chromedp.WaitReady("title", chromedp.After(func(_ context.Context, _ runtime.ExecutionContextID, _ ...*cdp.Node) error {
			<-gotResult

			return nil
		})),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		return "", fmt.Errorf("get body error: %w", err)
	}

	return body, nil
}

func (s *service) RenderPages(ctx context.Context, pages []*models.PageData, maxWorkers int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[0:], []chromedp.ExecAllocatorOption{
		chromedp.UserDataDir("./cache"),
		chromedp.Flag("new-window", false),
		chromedp.Flag("headless", s.prerenderConfig.Lookup.Headless),
		chromedp.UserAgent(s.prerenderConfig.UserAgent),
	}...)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	actxt, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// start the browser without a timeout
	if err := chromedp.Run(actxt); err != nil {
		panic(err)
	}

	sem := semaphore.NewWeighted(int64(maxWorkers))
	total := len(pages)

	for i, page := range pages {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)

			break
		}

		num := i

		go func(p *models.PageData) {
			defer sem.Release(1)

			if err := s.RenderPage(actxt, p, num, total); err != nil {
				log.Println(err)

				return
			}

			if p.Status != http.StatusOK {
				log.Printf("page http status is not 200. skip!")

				return
			}

			p.SuccessRender = true

			if err := s.r.SaveData(ctx, p); err != nil {
				log.Printf("save data error: %v", err)
			} else {
				p.SuccessStoring = true
			}

			// clear body to release memory usage
			p.Body = nil
		}(page)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		return fmt.Errorf("semafore acquire error: %w", err)
	}

	return nil
}

func (s *service) RenderPage(ctx context.Context, page *models.PageData, num, total int) error {
	if page == nil {
		return errors.ErrPageIsNil
	}

	page.Attempts++
	if page.Attempts == s.prerenderConfig.MaxAttempts {
		return fmt.Errorf("`%s` %w", page.URL.String(), errors.ErrMaxAttemptsExceed)
	}

	const logStatusFormat = "| %04d/%04d | %s | %d | %s"

	if err := s.GetPageBody(ctx, page); err != nil {
		log.Printf(logStatusFormat, num, total, "x", page.Attempts, fmt.Sprintf("%s | %s", page.URL.String(), err.Error()))

		// next attempt
		return s.RenderPage(ctx, page, num, total)
	}

	log.Printf(logStatusFormat, num, total, "v", page.Attempts, page.URL.String())

	return nil
}
