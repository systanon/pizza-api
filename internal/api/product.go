package api

import "time"

type GetProductQueryList struct {
	CategoryID *int64  `form:"categoryId"`
	Offset     *int    `form:"offset"`
	Limit      *int    `form:"limit"`
	Page       *int    `form:"page"`
	PerPage    *int    `form:"perPage"`
	Q          *string `form:"q"`
	SortOrder  string  `form:"sortOrder"`
}

type Product struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}
