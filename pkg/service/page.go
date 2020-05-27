package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	chromeRuntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/semaphore"

	"github.com/spacetab-io/prerender-go/pkg/models"
)

func (s service) GetPageBody(ctx context.Context, p *models.PageData) (err error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second) //nolint:gomnd
	defer cancel()

	newTabCtx, cancelNewTabCtx := chromedp.NewContext(timeoutCtx) // create new tab
	defer cancelNewTabCtx()

	var body string

	switch s.cfg.WaitFor {
	case models.WaitForConsole:
		body, err = s.renderBodyWithConsoleTrigger(newTabCtx, p)
	case models.WaitForElement:
		body, err = s.renderBodyWithElementTrigger(newTabCtx, p)
	case models.WaitForTime:
		body, err = s.renderBodyWithTimeTrigger(newTabCtx, p)
	default:
		err = errors.New("don't know what to wait")
	}

	if err != nil {
		return err
	}

	p.Body = []byte(body)
	p.ContentLength = len(body)
	p.Status = 200 //TODO убрать хардкод

	return nil
}

func (s service) renderBodyWithElementTrigger(ctx context.Context, p *models.PageData) (string, error) {
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.cfg.Viewport.Width, s.cfg.Viewport.Height, 1.0, false),
		chromedp.WaitVisible(s.cfg.Element.GetWaitElement()),
		chromedp.WaitVisible(s.cfg.Element.GetWaitElementAttr("ready")),
		chromedp.OuterHTML("html", &body),
	)

	return body, err
}
func (s service) renderBodyWithTimeTrigger(ctx context.Context, p *models.PageData) (string, error) {
	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.cfg.Viewport.Width, s.cfg.Viewport.Height, 1.0, false),
		chromedp.Sleep(s.cfg.SleepTime*time.Second),
		chromedp.OuterHTML("html", &body),
	)

	return body, err
}

func (s service) renderBodyWithConsoleTrigger(ctx context.Context, p *models.PageData) (string, error) {
	gotResult := make(chan bool, 1)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *chromeRuntime.EventConsoleAPICalled:
			if ev.Type == chromeRuntime.APITypeLog {
				for _, arg := range ev.Args {
					if string(arg.Value) == fmt.Sprintf(`"%s"`, s.cfg.ConsoleString) {
						gotResult <- true
					}
				}
			}
		case *chromeRuntime.EventExceptionThrown:
		}
	})

	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(p.URL.String()),
		emulation.SetDeviceMetricsOverride(s.cfg.Viewport.Width, s.cfg.Viewport.Height, 1.0, false),
		chromedp.WaitReady("title", chromedp.After(func(ctx context.Context, node ...*cdp.Node) error {
			<-gotResult
			return nil
		})),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		return "", err
	}

	return body, nil
}

func (s *service) RenderPages(pages []*models.PageData, maxWorkers int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[0:], []chromedp.ExecAllocatorOption{
		chromedp.UserDataDir("./cache"),
		chromedp.Flag("new-window", false),
		//chromedp.Flag("headless", false),
		chromedp.UserAgent(s.cfg.UserAgent),
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

	var wg sync.WaitGroup

	for i, page := range pages {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		page := page
		num := i

		go func() {
			defer sem.Release(1)

			if err := s.RenderPage(actxt, page, num); err != nil {
				log.Println(err)
				return
			}

			page.SuccessRender = true

			wg.Add(1)

			go func() {
				defer wg.Done()

				if err := s.r.SaveData(page); err != nil {
					log.Fatal(fmt.Errorf("save data error: %v", err))
				}
			}()
		}()
	}

	wg.Wait()

	return sem.Acquire(ctx, int64(maxWorkers))
}

func parseURI(uri string) (*url.URL, error) {
	if strings.Contains(uri, "?") {
		uri += "&"
	} else {
		uri += "?"
	}

	uri += "prerender=true"

	return url.Parse(uri)
}

func (s *service) RenderPage(ctx context.Context, page *models.PageData, num int) error {
	if page == nil {
		return errors.New("page data is nil")
	}

	page.Attempts++
	if page.Attempts == 5 { //nolint:gomnd
		return fmt.Errorf("render page `%s` attempts exceeded", page.URL.String())
	}

	const logStatusFormat = "| %04d | %s | %d | %s"

	err := s.GetPageBody(ctx, page)
	if err != nil {
		log.Printf(logStatusFormat, num, "x", page.Attempts, page.URL.String())

		// next attempt
		return s.RenderPage(ctx, page, num)
	}

	log.Printf(logStatusFormat, num, "v", page.Attempts, page.URL.String())

	return err
}
