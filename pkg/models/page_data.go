package models

import (
	"bytes"
	"compress/gzip"
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

func (d *PageData) MakeFileName(gzip bool) {
	page := d.URL.Path
	if page == "/" {
		page += "index"
	}

	page = strings.Trim(page, "/")

	d.FileName = page + ".html"

	if gzip {
		d.FileName += ".gzip"
	}
}

func (d *PageData) PrepareData(gzipFile bool) error {
	if !gzipFile {
		return nil
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(d.Body); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	d.Body = b.Bytes()

	return nil
}
