package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spacetab-io/commands-go"
	"github.com/spacetab-io/configuration-go/stage"
	"github.com/spacetab-io/prerender-go/configuration"
	"github.com/spacetab-io/prerender-go/pkg/log"
	"github.com/spacetab-io/prerender-go/pkg/repository"
	"github.com/spacetab-io/prerender-go/pkg/service"
	"github.com/spacetab-io/prerender-go/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "prerender",
		Short: "Prerender service",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run url parsing and prerend pages with storing html in storage",
		RunE:  run,
	}
)

func Execute() {
	envStage := stage.NewEnvStage("development")

	config, err := configuration.Init(envStage, os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatal().Err(err).Msg("config init error")
	}

	if err := log.Init(&config.Log, envStage.String(), config.Info.GetAlias(), config.Info.GetVersion()); err != nil {
		log.Fatal().Err(err).Msg("logs init fail")
	}

	rootCmd.AddCommand(commands.VersionCmd, runCmd)

	log.Info().Msg(config.Info.Summary())

	if err := rootCmd.ExecuteContext(setCtx(*config)); err != nil {
		os.Exit(commands.CmdFailureCode)
	}
}

func setCtx(config configuration.Config) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, commands.CommandContextObjectKeyConfig, &config)

	return ctx
}

func getConfigs(ctx context.Context) (
	*configuration.Config,
	error,
) {
	cfg, ok := ctx.Value(commands.CommandContextObjectKeyConfig).(*configuration.Config)
	if !ok {
		return nil, fmt.Errorf("%w: config (%s)", commands.ErrBadContextValue, commands.CommandContextObjectKeyConfig)
	}

	return cfg, nil
}

func initCfgAndService(cmd *cobra.Command) (*configuration.Config, service.Service, error) {
	cfg, err := getConfigs(cmd.Context())
	if err != nil {
		return nil, nil, utils.WrappedError("run", "getConfigs", err)
	}

	repo, err := repository.NewRepository(cfg.Storage)
	if err != nil {
		log.Error().Err(err).Msg("repo init error")

		return nil, nil, utils.WrappedError("initCfgAndService", "NewRepository", err)
	}

	return cfg, service.NewService(repo, cfg.Prerender, cfg.Storage), nil
}
