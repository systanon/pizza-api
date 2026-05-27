package model

import "time"

type Cart struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CartItem struct {
	ID        int64     `db:"id"`
	CartID    string    `db:"cart_id"`
	ProductID int64     `db:"product_id"`
	VariantID *int64    `db:"variant_id"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CartItemAddon struct {
	CartItemID int64 `db:"cart_item_id"`
	AddonID    int64 `db:"addon_id"`
}

type CartItemDetail struct {
	CartItem
	ProductName  string
	ProductImage string
	VariantName  *string
	VariantUnit  *string
	VariantPrice *int64
	Addons       []CartItemAddonDetail
}

type CartItemAddonDetail struct {
	AddonID   int64
	AddonName string
	Price     int64
}

type CartDetail struct {
	Cart
	Items []CartItemDetail
}
