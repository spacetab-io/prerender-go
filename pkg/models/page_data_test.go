package models

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageData_MakeFileName(t *testing.T) {
	type testCase struct {
		url      *url.URL
		postfix  string
		fileName string
	}

	tcs := []testCase{
		{url: &url.URL{Path: "/", RawQuery: ""}, fileName: "index"},
		{url: &url.URL{Path: "/page", RawQuery: ""}, fileName: "page"},
		{url: &url.URL{Path: "/page/deeper", RawQuery: ""}, fileName: "page/deeper"},
		{url: &url.URL{Path: "/page/deeper", RawQuery: ""}, postfix: "page", fileName: "page/deeper-page"},
		{url: &url.URL{Path: "/", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "index-bar=1-baz=2-foo=3"},
		{url: &url.URL{Path: "/page", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "page-bar=1-baz=2-foo=3"},
		{url: &url.URL{Path: "/page/deeper", RawQuery: "bar=1&baz=2&foo=3"}, fileName: "page/deeper-bar=1-baz=2-foo=3"},
		{url: &url.URL{Path: "/page/deeper", RawQuery: "bar=1&baz=2&foo=3"}, postfix: "page", fileName: "page/deeper-bar=1-baz=2-foo=3-page"},
	}

	for _, tc := range tcs {
		pd := &PageData{URL: tc.url}
		pd.MakeFileName(tc.postfix)
		assert.Equal(t, tc.fileName, pd.FileName)
	}
}
