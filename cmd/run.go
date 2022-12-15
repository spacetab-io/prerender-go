package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/log"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, _ []string) error {
	cfg, srv, err := initCfgAndService(cmd)
	if err != nil {
		return err
	}

	links, err := srv.GetLinksForRender()
	if err != nil {
		log.Error().Err(err).Msg("get links for render error")

		return err
	}

	timeStart := time.Now()

	log.Info().Int("links counts", len(links)).
		Str("lookup strategy", cfg.Prerender.Lookup.Type).
		Str("render wait strategy", cfg.Prerender.WaitFor).
		Msg("start rendering pages")

	pages, err := srv.PreparePages(links)
	if err != nil {
		log.Error().Err(err).Msg("prepare pages error")

		return fmt.Errorf("prepare pages error: %w", err)
	}

	maxWorkers := countMaxWorkers(cfg)

	if err := srv.RenderPages(cmd.Context(), pages, maxWorkers); err != nil {
		log.Error().Err(err).Msg("render pages error")

		return fmt.Errorf("render pages error: %w", err)
	}

	cmd.Print(srv.PrepareRenderReport(pages, time.Since(timeStart), maxWorkers))

	return nil
}

func countMaxWorkers(cfg *configuration.Config) int {
	numprocs := runtime.GOMAXPROCS(runtime.NumCPU())
	maxWorkers := 2 * numprocs //nolint:gomnd

	if cfg.Prerender.ConcurrentLimit == 0 {
		return numprocs
	}

	if cfg.Prerender.ConcurrentLimit > maxWorkers {
		return maxWorkers
	}

	return cfg.Prerender.ConcurrentLimit
}
