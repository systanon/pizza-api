package api

type CreateOrderRequest struct {
	CartID string `json:"cart_id" binding:"required"`
	Email  string `json:"email"   binding:"required,email"`
}

type OrderAddon struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type OrderItem struct {
	ID           int64        `json:"id"`
	ProductID    int64        `json:"product_id"`
	ProductName  string       `json:"product_name"`
	VariantID    *int64       `json:"variant_id,omitempty"`
	VariantName  string       `json:"variant_name,omitempty"`
	VariantPrice int64        `json:"variant_price"`
	Quantity     int          `json:"quantity"`
	ItemTotal    int64        `json:"item_total"`
	Addons       []OrderAddon `json:"addons"`
}

type OrderResponse struct {
	ID     int64       `json:"id"`
	CartID string      `json:"cart_id"`
	Email  string      `json:"email"`
	Status string      `json:"status"`
	Total  int64       `json:"total"`
	Items  []OrderItem `json:"items"`
}
