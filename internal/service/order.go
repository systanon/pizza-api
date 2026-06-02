package service

import (
	"context"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, cartID, email string, items []model.OrderItemInput, total int64) (*model.Order, error)
	GetOrder(ctx context.Context, orderID int64) (*model.OrderDetail, error)
}

type CartReader interface {
	GetCart(ctx context.Context, cartID string) (*model.CartDetail, error)
}

type OrderService struct {
	repo       OrderRepository
	cartReader CartReader
}

func NewOrderService(r OrderRepository, c CartReader) *OrderService {
	return &OrderService{repo: r, cartReader: c}
}

func (s *OrderService) CreateOrder(ctx context.Context, cartID, email string) (*model.OrderDetail, error) {
	cart, err := s.cartReader.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, apperror.ErrCartEmpty
	}

	total := calcTotal(cart.Items)
	items := buildOrderItems(cart.Items)

	order, err := s.repo.CreateOrder(ctx, cartID, email, items, total)
	if err != nil {
		return nil, err
	}

	return s.repo.GetOrder(ctx, order.ID)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int64) (*model.OrderDetail, error) {
	return s.repo.GetOrder(ctx, orderID)
}

func buildOrderItems(cartItems []model.CartItemDetail) []model.OrderItemInput {
	items := make([]model.OrderItemInput, 0, len(cartItems))
	for _, ci := range cartItems {
		variantName := ""
		if ci.VariantName != nil {
			variantName = *ci.VariantName
		}
		variantPrice := int64(0)
		if ci.VariantPrice != nil {
			variantPrice = *ci.VariantPrice
		}

		addonTotal := int64(0)
		addons := make([]model.OrderAddonInput, 0, len(ci.Addons))
		for _, a := range ci.Addons {
			addonTotal += a.Price
			addons = append(addons, model.OrderAddonInput{
				AddonID:   a.AddonID,
				AddonName: a.AddonName,
				Price:     a.Price,
			})
		}
		itemTotal := (variantPrice + addonTotal) * int64(ci.Quantity)

		items = append(items, model.OrderItemInput{
			ProductID:    ci.ProductID,
			ProductName:  ci.ProductName,
			VariantID:    ci.CartItem.VariantID,
			VariantName:  variantName,
			VariantPrice: variantPrice,
			Quantity:     ci.Quantity,
			ItemTotal:    itemTotal,
			Addons:       addons,
		})
	}
	return items
}
