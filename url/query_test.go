package url

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareQueryParams(t *testing.T) {
	type testCase struct {
		inputQuery   string
		inputParams  []string
		outputString string
	}

	tcs := []testCase{
		{inputParams: []string{"baz", "var", "foo"}, inputQuery: "var=1&foo=2&baz=3", outputString: "baz=3&foo=2&var=1"},
		{inputParams: []string{"baz", "var", "foo"}, inputQuery: "foo=2&baz=3", outputString: "baz=3&foo=2"},
		{inputParams: []string{"var", "foo"}, inputQuery: "var=1&baz=3&foo=2", outputString: "foo=2&var=1"},
	}

	for _, tc := range tcs {
		uri := &url.URL{RawQuery: tc.inputQuery}
		PrepareSortedQueryParams(uri, tc.inputParams)

		assert.Equal(t, tc.outputString, uri.RawQuery)
	}
}
