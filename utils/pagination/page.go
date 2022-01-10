package pagination

import (
	"net/url"
	"strconv"
)

type Page struct {
	Page uint64 `json:"page,omitempty"`
	Size uint64 `json:"size,omitempty"`
}

func (p *Page) GetLimit() uint64 {
	return p.Size
}

func (p *Page) GetOffset() uint64 {
	return p.Size * (p.Page - 1)
}

func CreatePage(params url.Values) *Page {
	page, err := strconv.ParseUint(params.Get("page"), 10, 64)
	size, err2 := strconv.ParseUint(params.Get("size"), 10, 64)
	if page == 0 || size == 0 || err != nil || err2 != nil {
		return &Page{Page: 1, Size: 10}
	}
	return &Page{Page: page, Size: size}

}
