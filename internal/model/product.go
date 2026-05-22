package model

import "time"

type Product struct {
	ID          int64            `db:"id"`
	CategoryID  int64            `db:"category_id"`
	Name        string           `db:"name"`
	Description string           `db:"description"`
	ImageURL    string           `db:"image_url"`
	Metadata    *ProductMetadata `db:"metadata" json:"metadata"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

type ProductVariantPrice struct {
	ProductID int64     `db:"product_id"`
	VariantID int64     `db:"variant_id"`
	Price     int64     `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

type ProductPromo struct {
	IsActive      bool       `json:"is_active"`
	DiscountType  string     `json:"discount_type,omitempty"`
	DiscountValue int64      `json:"discount_value,omitempty"`
	PromoText     string     `json:"promo_text,omitempty"`
	ValidUntil    *time.Time `json:"valid_until,omitempty"`
}

type ProductMetadata struct {
	Calories int           `json:"calories,omitempty"`
	IsVeggie bool          `json:"is_veggie,omitempty"`
	IsSpicy  bool          `json:"is_spicy,omitempty"`
	Tags     []string      `json:"tags,omitempty"`
	Promo    *ProductPromo `json:"promo,omitempty"`
}

type ProductDetailRepositoryDTO struct {
	Product  Product
	Variants []VariantRepositoryPriceDTO
}

type ProductDTO struct {
	Product  Product
	Variants []VariantPriceDTO
}
