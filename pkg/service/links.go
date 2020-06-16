package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/yterajima/go-sitemap"

	"github.com/spacetab-io/prerender-go/pkg/models"
	prerenderUrl "github.com/spacetab-io/prerender-go/url"
)

func (s *service) PreparePages(links []string) ([]*models.PageData, error) {
	pages := make([]*models.PageData, 0)

	for _, link := range links {
		uri, err := url.Parse(link)
		if err != nil {
			return nil, fmt.Errorf("parse link url error: %v", err)
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

	return nil, errors.New("lookup type is wrong or not set")
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
		return nil, errors.New("base_url is not set in config")
	}

	var links = make([]string, 0)

	for _, link := range s.prerenderConfig.Lookup.PageURLs {
		if strings.Contains(link, "https://") {
			return nil, errors.New("link contains base url")
		}

		links = append(links, s.prerenderConfig.Lookup.BaseURL+link)
	}

	return links, nil
}

func (s *service) GetUrlsFromSitemaps() ([]string, error) {
	links := make([]string, 0)

	for _, url := range s.prerenderConfig.Lookup.SitemapURLs {
		smap, err := sitemap.Get(url, nil)
		if err != nil {
			return nil, err
		}

		for _, URL := range smap.URL {
			if !IsInSlice(links, URL.Loc) {
				links = append(links, URL.Loc)
			}
		}
	}

	return links, nil
}
