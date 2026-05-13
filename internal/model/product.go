package model

import "time"

type Product struct {
	ID          int64     `json:"id" db:"id"`
	CategoryID  int64     `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       int64     `json:"price" db:"price"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ProductFilter struct {
	CategoryID *int64
	Q          *string
	SortOrder  string
	Limit      int
	Offset     int
}

type ProductQuery struct {
	CategoryID *int64  `form:"categoryId"`
	Offset     *int    `form:"offset"`
	Limit      *int    `form:"limit"`
	Page       *int    `form:"page"`
	PerPage    *int    `form:"perPage"`
	Q          *string `form:"q"`
	SortOrder  string  `form:"sortOrder"`
}
