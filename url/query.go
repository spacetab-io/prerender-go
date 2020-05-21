package url

import (
	"net/url"
	"sort"
)

//PrepareSortedQueryParams Extract only query params that are needed to preserve and sort them alphabetically in result query string
func PrepareSortedQueryParams(uri *url.URL, queryParamsToSave []string) {
	query := uri.Query()
	queryParams := make(url.Values)
	keys := make([]string, 0)

	for _, key := range queryParamsToSave {
		val := query.Get(key)
		if val != "" {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	for _, k := range keys {
		queryParams.Add(k, query.Get(k))
	}

	uri.RawQuery = queryParams.Encode()
}
