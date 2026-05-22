package api

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
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type ProductDetail struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Variants    []Variant `json:"variants,omitempty"`
}
