package service

import (
	"context"
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

type CartService struct {
	repo CartRepository
}

func NewCartService(r CartRepository) *CartService {
	return &CartService{repo: r}
}

func (s *CartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	return s.repo.CreateCart(ctx)
}

func (s *CartService) GetCart(ctx context.Context, cartID string) (*model.CartDetail, error) {
	return s.repo.GetCart(ctx, cartID)
}

func (s *CartService) AddItem(ctx context.Context, cartID string, productID int64, variantID *int64, quantity int, addonIDs []int64) (*model.CartItem, error) {
	if addonIDs == nil {
		addonIDs = []int64{}
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
