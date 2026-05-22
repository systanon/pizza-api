package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	if total == 0 {
		return products, 0, nil
	}

	productQuery := fmt.Sprintf(
		"SELECT id, category_id, name, description, image_url, created_at, updated_at %s ORDER BY created_at %s LIMIT $%d OFFSET $%d",
		baseQuery,
		f.SortOrder,
		argIndex,
		argIndex+1,
	)

	args = append(args, f.Limit, f.Offset)

	err = r.db.SelectContext(ctx, &products, productQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, id int64) (*model.ProductDetailRepositoryDTO, error) {

	query := `
		SELECT 
			p.id, p.category_id, p.name, p.description, p.image_url, p.created_at,
			v.id AS variant_id, v.name AS variant_name, v.value, v.unit, v.created_at AS variant_created_at, v.updated_at AS variant_updated_at,
			pvp.price
		FROM products p
		LEFT JOIN product_variant_prices pvp ON p.id = pvp.product_id
		LEFT JOIN variants v ON v.id = pvp.variant_id
		WHERE p.id = $1`

	type flatRow struct {
		model.Product
		VariantID        sql.NullInt64  `db:"variant_id"`
		VariantName      sql.NullString `db:"variant_name"`
		VariantValue     sql.NullInt64  `db:"value"`
		VariantUnit      sql.NullString `db:"unit"`
		Price            sql.NullInt64  `db:"price"`
		VariantCreatedAt sql.NullTime   `db:"variant_created_at"`
		VariantUpdatedAt sql.NullTime   `db:"variant_updated_at"`
	}

	var rows []flatRow
	err := r.db.SelectContext(ctx, &rows, query, id)

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, apperror.ErrNotFound
	}

	result := &model.ProductDetailRepositoryDTO{
		Product:  rows[0].Product,
		Variants: []model.VariantRepositoryPriceDTO{},
	}

	for _, row := range rows {
		if row.VariantID.Valid {

			var createdAt, updatedAt time.Time
			if row.VariantCreatedAt.Valid {
				createdAt = row.VariantCreatedAt.Time
			}
			if row.VariantUpdatedAt.Valid {
				updatedAt = row.VariantUpdatedAt.Time
			}
			result.Variants = append(result.Variants, model.VariantRepositoryPriceDTO{
				ID:        row.VariantID.Int64,
				Name:      row.VariantName.String,
				Value:     row.VariantValue.Int64,
				Unit:      row.VariantUnit.String,
				Price:     row.Price.Int64,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			})
		}
	}

	return result, nil
}
