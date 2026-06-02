package apperror

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("entity not found")
	ErrVariantRequired = errors.New("variant is required for this product")
	ErrVariantInvalid  = errors.New("variant does not belong to this product")
	ErrNoPrice         = errors.New("product has no price defined")
	ErrCartEmpty        = errors.New("cart is empty")
	ErrOrderNotPending  = errors.New("order is not in pending status")
)
