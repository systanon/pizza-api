package service

import (
	"context"
	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"
)

type CartRepository interface {
	CreateCart(ctx context.Context) (*model.Cart, error)
	GetCart(ctx context.Context, cartID string) (*model.CartDetail, error)
	AddItem(ctx context.Context, cartID string, productID int64, variantID *int64, quantity int, addonIDs []int64) (*model.CartItem, error)
	UpdateItem(ctx context.Context, cartID string, itemID int64, quantity int) error
	RemoveItem(ctx context.Context, cartID string, itemID int64) error
	ClearCart(ctx context.Context, cartID string) error
}

type ProductChecker interface {
	HasVariants(ctx context.Context, productID int64) (bool, error)
	IsValidVariant(ctx context.Context, productID int64, variantID int64) (bool, error)
}

type CartService struct {
	repo           CartRepository
	productChecker ProductChecker
}

func NewCartService(r CartRepository, p ProductChecker) *CartService {
	return &CartService{repo: r, productChecker: p}
}

func (s *CartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	return s.repo.CreateCart(ctx)
}

func (s *CartService) GetCart(ctx context.Context, cartID string) (*model.CartDetail, error) {
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}
	cart.Total = calcTotal(cart.Items)
	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, cartID string, productID int64, variantID *int64, quantity int, addonIDs []int64) (*model.CartItem, error) {
	if addonIDs == nil {
		addonIDs = []int64{}
	}

	hasVariants, err := s.productChecker.HasVariants(ctx, productID)
	if err != nil {
		return nil, err
	}

	if !hasVariants {
		return nil, apperror.ErrNoPrice
	}

	if variantID == nil {
		return nil, apperror.ErrVariantRequired
	}

	valid, err := s.productChecker.IsValidVariant(ctx, productID, *variantID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, apperror.ErrVariantInvalid
	}

	return s.repo.AddItem(ctx, cartID, productID, variantID, quantity, addonIDs)
}

func (s *CartService) UpdateItem(ctx context.Context, cartID string, itemID int64, quantity int) error {
	return s.repo.UpdateItem(ctx, cartID, itemID, quantity)
}

func (s *CartService) RemoveItem(ctx context.Context, cartID string, itemID int64) error {
	return s.repo.RemoveItem(ctx, cartID, itemID)
}

func (s *CartService) ClearCart(ctx context.Context, cartID string) error {
	return s.repo.ClearCart(ctx, cartID)
}

func calcTotal(items []model.CartItemDetail) int64 {
	var total int64
	for _, item := range items {
		if item.VariantPrice != nil {
			total += *item.VariantPrice * int64(item.Quantity)
		}
		for _, addon := range item.Addons {
			total += addon.Price * int64(item.Quantity)
		}
	}
	return total
}
