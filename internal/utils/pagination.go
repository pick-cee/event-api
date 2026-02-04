package utils

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// Generic paginated response
type PaginatedResponse[T any] struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	Data       []T   `json:"data"`
}

// Pagination params
type PaginationParams struct {
	Page  int
	Limit int
}

func GetPaginationParams(r *http.Request) PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

func Paginate(params PaginationParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.Limit
		return db.Offset(offset).Limit(params.Limit)
	}
}

func NewPaginationResponse[T any](data []T, total int64, params PaginationParams) PaginatedResponse[T] {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return PaginatedResponse[T]{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: totalPages,
		Data:       data,
	}
}
