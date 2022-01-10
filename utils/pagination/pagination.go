package pagination

type Pagination struct {
	Page
	Count uint64      `json:"count,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func CreatePagination(page Page, data interface{}, count uint64) *Pagination {
	return &Pagination{Page: page, Data: data, Count: count}
}
