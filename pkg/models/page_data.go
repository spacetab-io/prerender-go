package models

import (
	"net/url"
	"strings"
)

type PageData struct {
	URL           *url.URL
	Status        int
	ContentLength int
	Body          []byte
	FileName      string
	Attempts      int
	SuccessRender bool
}

func (d *PageData) MakeFileName() {
	page := d.URL.Path
	if page == "/" {
		page += "index"
	}

	page = strings.Trim(page, "/")

	q := ""
	if d.URL.RawQuery != "" {
		q += "-" + strings.ReplaceAll(d.URL.RawQuery, "&", "-")
	}

	d.FileName = page + q
}
