package repository

import (
	"context"

	"pizza-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetCategoryList(ctx context.Context) ([]model.Category, error) {
	categories := []model.Category{}
	err := r.db.SelectContext(ctx, &categories,
		"SELECT id, name, slug, created_at FROM categories ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	return categories, nil
}
