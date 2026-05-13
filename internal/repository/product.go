package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetProductList(ctx context.Context, f model.ProductFilter) ([]model.Product, int, error) {
	products := []model.Product{}
	var total int

	baseQuery := "FROM products WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if f.CategoryID != nil {
		baseQuery += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *f.CategoryID)
		argIndex++
	}

	if f.Q != nil {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+*f.Q+"%", "%"+*f.Q+"%")
		argIndex += 2
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	todoQuery := fmt.Sprintf(
		"SELECT id, category_id, name, description, price, image_url, created_at %s ORDER BY created_at %s LIMIT $%d OFFSET $%d",
		baseQuery,
		f.SortOrder,
		argIndex,
		argIndex+1,
	)

	args = append(args, f.Limit, f.Offset)

	err = r.db.SelectContext(ctx, &products, todoQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {

	product := model.Product{}
	err := r.db.GetContext(ctx, &product, "SELECT id, category_id, name, description, price, image_url, created_at FROM products WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return &product, nil
}
