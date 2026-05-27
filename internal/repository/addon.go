package repository

import (
	"context"

	"pizza-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type AddonRepo struct {
	db *sqlx.DB
}

func NewAddonRepo(db *sqlx.DB) *AddonRepo {
	return &AddonRepo{db: db}
}

func (r *AddonRepo) GetAddonList(ctx context.Context) ([]model.Addon, error) {
	addons := []model.Addon{}
	err := r.db.SelectContext(ctx, &addons,
		"SELECT id, name, price, created_at, updated_at FROM addons ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	return addons, nil
}

func (r *AddonRepo) GetAddonListByCategoryID(ctx context.Context, categoryID int64) ([]model.AddonRow, error) {
	addons := []model.AddonRow{}

	query := `
		SELECT a.id, a.name, a.price, a.created_at, a.updated_at
		FROM addons a
		INNER JOIN category_addons ca ON a.id = ca.addon_id
		WHERE ca.category_id = $1
		ORDER BY a.name ASC`

	err := r.db.SelectContext(ctx, &addons, query, categoryID)
	if err != nil {
		return nil, err
	}

	return addons, nil
}
