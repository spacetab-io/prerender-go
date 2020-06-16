package models

import (
	"net/url"
	"strings"
)

type PageData struct {
	URL            *url.URL
	Status         int
	ContentLength  int
	Body           []byte
	FileName       string
	Attempts       int
	SuccessRender  bool
	SuccessStoring bool
}

func (d *PageData) MakeFileName(postfix string) {
	page := d.URL.Path
	if page == "/" {
		page += "index"
	}

	page = strings.Trim(page, "/")

	fileNameElems := make([]string, 0)

	fileNameElems = append(fileNameElems, page)
	if d.URL.RawQuery != "" {
		fileNameElems = append(fileNameElems, strings.ReplaceAll(d.URL.RawQuery, "&", "-"))
	}

	if postfix != "" {
		fileNameElems = append(fileNameElems, postfix)
	}

	d.FileName = strings.Join(fileNameElems, "-")
}
