package service

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spacetab-io/prerender-go/pkg/errors"
	"github.com/spacetab-io/prerender-go/pkg/log"
	"github.com/spacetab-io/prerender-go/pkg/models"
	prerenderUrl "github.com/spacetab-io/prerender-go/url"
	"github.com/yterajima/go-sitemap"
)

func (s *service) PreparePages(links []string) ([]*models.PageData, error) {
	pages := make([]*models.PageData, 0)

	for _, link := range links {
		uri, err := url.Parse(link)
		if err != nil {
			return nil, fmt.Errorf("parse link url error: %w", err)
		}

		prerenderUrl.PrepareSortedQueryParams(uri, s.prerenderConfig.Lookup.ParamsToSave)

		page := &models.PageData{URL: uri, Attempts: 0}
		page.MakeFileName(s.prerenderConfig.FilePostfix)
		appendPrerenderQueryParam(page)
		pages = append(pages, page)
	}

	return pages, nil
}

func appendPrerenderQueryParam(page *models.PageData) {
	if len(page.URL.RawQuery) > 0 {
		page.URL.RawQuery += "&"
	}

	page.URL.RawQuery += "prerender=true"
}

func (s *service) GetLinksForRender() ([]string, error) {
	switch s.prerenderConfig.Lookup.Type {
	case models.LookupTypeSitemaps:
		return s.GetUrlsFromSitemaps()
	case models.LookupTypeURLs:
		return s.GetUrlsFromLinksList()
	case models.LookupTypeAll:
		links, err := s.GetUrlsFromSitemaps()
		if err != nil {
			return nil, err
		}

		configLinks, err := s.GetUrlsFromLinksList()
		if err != nil {
			return nil, err
		}

		for _, configLink := range configLinks {
			if !IsInSlice(links, configLink) {
				links = append(links, configLink)
			}
		}

		return links, nil
	}

	return nil, errors.ErrWrongLookupType
}

func IsInSlice(links []string, link string) bool {
	for _, l := range links {
		if l == link {
			return true
		}
	}

	return false
}

func (s *service) GetUrlsFromLinksList() ([]string, error) {
	if s.prerenderConfig.Lookup.BaseURL == "" {
		return nil, errors.ErrNoBaseURL
	}

	links := make([]string, 0)

	for _, link := range s.prerenderConfig.Lookup.PageURLs {
		if strings.Contains(link, "https://") {
			return nil, errors.ErrLinkHasBaseURL
		}

		links = append(links, s.prerenderConfig.Lookup.BaseURL+link)
	}

	return links, nil
}

func (s *service) GetUrlsFromSitemaps() ([]string, error) {
	links := make([]string, 0)

	for _, sitemapsURL := range s.prerenderConfig.Lookup.SitemapURLs {
		smap, err := sitemap.Get(sitemapsURL, nil)
		if err != nil {
			return nil, err
		}

		for _, sitemapURL := range smap.URL {
			if sitemapURL.LastMod == "" {
				log.Warn().Str("sitemap url", sitemapURL.Loc).Msg("last modified is empty")

				continue
			}

			lm, err := lastModifiedFrom(sitemapURL.LastMod)
			if err != nil {
				return nil, err
			}

			if s.lastRenderedAt.After(lm) {
				continue
			}

			if !IsInSlice(links, sitemapURL.Loc) {
				links = append(links, sitemapURL.Loc)
			}
		}
	}

	return links, nil
}

func lastModifiedFrom(lastMod string) (time.Time, error) {
	// 2021-09-21T07:31:56+00:00
	t, err := time.Parse(time.RFC3339, lastMod)
	if err != nil {
		return time.Time{}, fmt.Errorf("lastModifiedFrom parse error: %w", err)
	}

	return t, nil
}
