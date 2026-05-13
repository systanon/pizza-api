package service

import (
	"context"
	"pizza-backend/internal/model"
)

type CategoryRepository interface {
	GetCategoryList(ctx context.Context) ([]model.Category, error)
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(r CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetCategoryList(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetCategoryList(ctx)
}
