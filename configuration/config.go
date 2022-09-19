package configuration

import (
	"fmt"

	config "github.com/spacetab-io/configuration-go"
	"github.com/spacetab-io/configuration-go/stage"
	cfgstructs "github.com/spacetab-io/configuration-structs-go/v2"
	"github.com/spacetab-io/prerender-go/pkg/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Info      cfgstructs.ApplicationInfo `yaml:"info"`
	Log       cfgstructs.Logs            `yaml:"logs"`
	Storage   StorageConfig              `yaml:"storage"`
	Prerender PrerenderConfig            `yaml:"prerender"`
}

func Init(s stage.Interface, path string) (*Config, error) {
	configBytes, err := config.Read(s, path, false)
	if err != nil {
		log.Error().Err(err).Msg("read error")

		return nil, fmt.Errorf("init config error: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		log.Error().Err(err).Bytes("error cfg", configBytes).Msg("unmarshal error for bytes")

		return nil, fmt.Errorf("config unmarshal error: %w", err)
	}

	return &cfg, nil
}
