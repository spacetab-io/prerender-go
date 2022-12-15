package models_test

import (
	"net/url"
	"testing"

	"github.com/spacetab-io/prerender-go/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestPageData_MakeFileName(t *testing.T) {
	type testCase struct {
		name     string
		url      *url.URL
		postfix  string
		fileName string
	}

	tcs := []testCase{
		{name: "1", url: &url.URL{Path: "/", RawQuery: ""}, fileName: "index"},
		{name: "2", url: &url.URL{Path: "/page", RawQuery: ""}, fileName: "page"},
		{name: "3", url: &url.URL{Path: "/page/deeper", RawQuery: ""}, fileName: "page/deeper"},
		{name: "4", url: &url.URL{Path: "/page/deeper", RawQuery: ""}, postfix: "page", fileName: "page/deeper-page"},
		{name: "5", url: &url.URL{Path: "/", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "index-bar=1-baz=2-foo=3"},
		{name: "6", url: &url.URL{Path: "/page", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "page-bar=1-baz=2-foo=3"},
		{name: "7", url: &url.URL{Path: "/page/deeper", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "page/deeper-bar=1-baz=2-foo=3"},
		{name: "8", url: &url.URL{Path: "/page/deeper", RawQuery: "bar=1&baz=2&foo=3"}, postfix: "page", fileName: "page/deeper-bar=1-baz=2-foo=3-page"},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pd := &models.PageData{URL: tc.url}
			pd.MakeFileName(tc.postfix)
			assert.Equal(t, tc.fileName, pd.FileName)
		})
	}
}
