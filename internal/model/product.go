package model

import "time"

type Product struct {
	ID          int64     `db:"id"`
	CategoryID  int64     `db:"category_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       int64     `db:"price"`
	ImageURL    string    `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}

type ProductFilter struct {
	CategoryID *int64
	Q          *string
	SortOrder  string
	Limit      int
	Offset     int
}

type ProductParams struct {
	CategoryID *int64
	Offset     *int
	Limit      *int
	Page       *int
	PerPage    *int
	Q          *string
	SortOrder  string
}
