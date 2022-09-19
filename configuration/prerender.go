package configuration

import (
	"fmt"
	"strings"
	"time"

	"github.com/spacetab-io/prerender-go/pkg/models"
)

type ElementConfig struct {
	Type      string `yaml:"type"`
	ID        string `yaml:"id"`
	Class     string `yaml:"class"`
	Attribute struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"attribute"`
}

type lookupConfig struct {
	Headless     bool     `yaml:"headless"`
	Type         string   `yaml:"type"`
	SitemapURLs  []string `yaml:"sitemaps"`
	PageURLs     []string `yaml:"urls"`
	BaseURL      string   `yaml:"base_url"`
	ParamsToSave []string `yaml:"get_params_to_save"`
}

func (c lookupConfig) GetSourceURL() string {
	switch c.Type {
	case models.LookupTypeSitemaps:
		return strings.Join(c.SitemapURLs, ",")
	case models.LookupTypeURLs:
		return strings.Join(c.PageURLs, ", ")
	}

	return ""
}

type viewportConfig struct {
	Width  int64 `yaml:"width"`
	Height int64 `yaml:"height"`
}

type PrerenderConfig struct {
	UserAgent       string         `yaml:"user_agent"`
	FilePostfix     string         `yaml:"file_postfix"`
	ConcurrentLimit int            `yaml:"concurrent_limit"`
	Lookup          lookupConfig   `yaml:"lookup"`
	WaitFor         string         `yaml:"wait_for"`
	ConsoleString   string         `yaml:"console_string"`
	SleepTime       time.Duration  `yaml:"sleep_time"`
	Element         ElementConfig  `yaml:"element"`
	Viewport        viewportConfig `yaml:"viewport"`
}

func (ec ElementConfig) GetWaitElement() string {
	elem := ec.Type

	if ec.ID != "" {
		elem = "#" + ec.ID
	}

	if ec.Class != "" && ec.ID == "" {
		elem = "." + ec.Class
	}

	return elem
}

func (ec ElementConfig) GetWaitElementAttr(attrValue string) string {
	elem := ec.GetWaitElement()

	if ec.Attribute.Name != "" && ec.Attribute.Value != "" {
		elem += fmt.Sprintf("[%s=%s]", ec.Attribute.Name, attrValue)
	}

	return elem
}
