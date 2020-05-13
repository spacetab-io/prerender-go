package main

import (
	"log"
	"runtime"
	"time"

	cfg "github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/service"
	"github.com/spacetab-io/prerender-go/pkg/storage"
)

func main() {
	if err := cfg.Init(); err != nil {
		log.Fatalf("config reading error: %+v", err)
	}

	st, err := storage.NewStorage(cfg.Config.Storage)

	if err != nil {
		log.Fatal(err)
	}

	srv := service.NewService(st, cfg.Config.Prerender)

	links, err := srv.GetLinksForRender()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("links counts %d\n", len(links))

	pages, err := srv.PreparePages(links)
	if err != nil {
		log.Fatal(err)
	}

	timeStart := time.Now()
	maxWorkers := runtime.GOMAXPROCS(runtime.NumCPU()) * cfg.Config.Prerender.ChromeTabsMultiplier

	if err := srv.RenderPages(pages, maxWorkers); err != nil {
		log.Fatal(err)
	}

	timeEnd := time.Now()

	srv.PrepareRenderReport(pages, timeEnd.Sub(timeStart), maxWorkers)
}
