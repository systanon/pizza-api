package model

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        int64       `db:"id"`
	CartID    string      `db:"cart_id"`
	Email     string      `db:"email"`
	Status    OrderStatus `db:"status"`
	Total     int64       `db:"total"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

type OrderItem struct {
	ID           int64  `db:"id"`
	OrderID      int64  `db:"order_id"`
	ProductID    int64  `db:"product_id"`
	ProductName  string `db:"product_name"`
	VariantID    *int64 `db:"variant_id"`
	VariantName  string `db:"variant_name"`
	VariantPrice int64  `db:"variant_price"`
	Quantity     int    `db:"quantity"`
	ItemTotal    int64  `db:"item_total"`
}

type OrderItemAddon struct {
	OrderItemID int64  `db:"order_item_id"`
	AddonID     int64  `db:"addon_id"`
	AddonName   string `db:"addon_name"`
	Price       int64  `db:"price"`
}

type OrderItemDetail struct {
	OrderItem
	Addons []OrderItemAddon
}

type OrderDetail struct {
	Order
	Items []OrderItemDetail
}

// OrderItemInput is a pre-computed snapshot passed from service to repository.
type OrderItemInput struct {
	ProductID    int64
	ProductName  string
	VariantID    *int64
	VariantName  string
	VariantPrice int64
	Quantity     int
	ItemTotal    int64
	Addons       []OrderItemAddon
}
