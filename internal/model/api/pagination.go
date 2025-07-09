package model

import (
	"net/http"
	"strconv"
)

type PaginationResponse[T any] struct {
	Data       []T         `json:"data"`
	Pagination interface{} `json:"pagination"`
}

type Pagination struct {
	Page       uint `json:"page"`
	PageSize   uint `json:"page_size"`
	TotalRows  uint `json:"total_rows"`
	TotalPages uint `json:"total_pages"`
}

func NewPaginationFromRequest(r *http.Request) *Pagination {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 50
	}

	return &Pagination{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
}

func (p *Pagination) Limit() uint {
	return p.PageSize
}

func (p *Pagination) Offset() uint {
	return (p.Page - 1) * p.PageSize
}
