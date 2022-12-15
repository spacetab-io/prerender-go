package url_test

import (
	"net/url"
	"testing"

	url2 "github.com/spacetab-io/prerender-go/url"
	"github.com/stretchr/testify/assert"
)

func TestPrepareQueryParams(t *testing.T) {
	type testCase struct {
		name         string
		inputQuery   string
		inputParams  []string
		outputString string
	}

	tcs := []testCase{
		{name: "1", inputParams: []string{"baz", "var", "foo"}, inputQuery: "var=1&foo=2&baz=3", outputString: "baz=3&foo=2&var=1"},
		{name: "2", inputParams: []string{"baz", "var", "foo"}, inputQuery: "foo=2&baz=3", outputString: "baz=3&foo=2"},
		{name: "3", inputParams: []string{"var", "foo"}, inputQuery: "var=1&baz=3&foo=2", outputString: "foo=2&var=1"},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uri := &url.URL{RawQuery: tc.inputQuery}
			url2.PrepareSortedQueryParams(uri, tc.inputParams)

			assert.Equal(t, tc.outputString, uri.RawQuery)
		})
	}
}
