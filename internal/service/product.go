package service

import (
	"context"
	"pizza-backend/internal/model"
)

type ProductRepository interface {
	GetProductList(ctx context.Context, f model.ProductFilter) ([]model.Product, int, error)
	GetProductByID(ctx context.Context, id int64) (*model.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) GetProductList(ctx context.Context, q model.ProductParams) ([]model.Product, int, int, error) {
	const defaultLimit = 10
	const defaultSortOrder = "DESC"

	offset := 0
	limit := defaultLimit

	sortOrder := q.SortOrder
	if q.SortOrder != "ASC" && q.SortOrder != "DESC" {
		sortOrder = defaultSortOrder
	}

	if q.Page != nil && *q.Page > 0 && q.PerPage != nil && *q.PerPage > 0 {
		page := *q.Page
		perPage := *q.PerPage
		offset = (page - 1) * perPage
		limit = perPage
	} else if q.Offset != nil || q.Limit != nil {

		if q.Offset != nil {
			offset = *q.Offset
		}

		if q.Limit != nil {
			limit = *q.Limit
		}
	}

	filter := model.ProductFilter{
		CategoryID: q.CategoryID,
		Q:          q.Q,
		SortOrder:  sortOrder,
		Limit:      limit,
		Offset:     offset,
	}

	products, count, err := s.repo.GetProductList(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	pages := (count + limit - 1) / limit

	return products, count, pages, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}
