package api

type CreateCartItemRequest struct {
	ProductID int64   `json:"product_id" binding:"required"`
	VariantID *int64  `json:"variant_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	AddonIDs  []int64 `json:"addon_ids"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartAddon struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type CartItem struct {
	ID           int64       `json:"id"`
	ProductID    int64       `json:"product_id"`
	ProductName  string      `json:"product_name"`
	ProductImage string      `json:"product_image"`
	VariantID    *int64      `json:"variant_id,omitempty"`
	VariantName  *string     `json:"variant_name,omitempty"`
	VariantUnit  *string     `json:"variant_unit,omitempty"`
	VariantPrice *int64      `json:"variant_price,omitempty"`
	Quantity     int         `json:"quantity"`
	Addons       []CartAddon `json:"addons"`
}

type CartResponse struct {
	ID    string     `json:"id"`
	Items []CartItem `json:"items"`
	Total int64      `json:"total"`
}
