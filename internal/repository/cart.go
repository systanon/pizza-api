package repository

import (
	"context"
	"database/sql"
	"strings"

	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type CartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) CreateCart(ctx context.Context) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO carts DEFAULT VALUES RETURNING id, created_at, updated_at`,
	).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepo) GetCart(ctx context.Context, cartID string) (*model.CartDetail, error) {
	var cart model.Cart
	err := r.db.GetContext(ctx, &cart,
		`SELECT id, created_at, updated_at FROM carts WHERE id = $1`, cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	type itemRow struct {
		ID           int64          `db:"id"`
		CartID       string         `db:"cart_id"`
		ProductID    int64          `db:"product_id"`
		VariantID    sql.NullInt64  `db:"variant_id"`
		Quantity     int            `db:"quantity"`
		ProductName  string         `db:"product_name"`
		ProductImage string         `db:"product_image"`
		VariantName  sql.NullString `db:"variant_name"`
		VariantUnit  sql.NullString `db:"variant_unit"`
		VariantPrice sql.NullInt64  `db:"variant_price"`
	}

	var rows []itemRow
	err = r.db.SelectContext(ctx, &rows, `
		SELECT
			ci.id,
			ci.cart_id,
			ci.product_id,
			ci.variant_id,
			ci.quantity,
			p.name        AS product_name,
			p.image_url   AS product_image,
			v.name        AS variant_name,
			v.unit        AS variant_unit,
			pvp.price     AS variant_price
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		LEFT JOIN variants v ON v.id = ci.variant_id
		LEFT JOIN product_variant_prices pvp
			ON pvp.product_id = ci.product_id AND pvp.variant_id = ci.variant_id
		WHERE ci.cart_id = $1
		ORDER BY ci.id`, cartID)
	if err != nil {
		return nil, err
	}

	itemIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		itemIDs = append(itemIDs, row.ID)
	}

	addonsByItem := map[int64][]model.CartItemAddonDetail{}
	if len(itemIDs) > 0 {
		type addonRow struct {
			CartItemID int64  `db:"cart_item_id"`
			AddonID    int64  `db:"addon_id"`
			AddonName  string `db:"addon_name"`
			Price      int64  `db:"price"`
		}
		var addonRows []addonRow
		query, args, err := sqlx.In(`
			SELECT cia.cart_item_id, cia.addon_id, a.name AS addon_name, a.price
			FROM cart_item_addons cia
			JOIN addons a ON a.id = cia.addon_id
			WHERE cia.cart_item_id IN (?)`, itemIDs)
		if err != nil {
			return nil, err
		}
		query = r.db.Rebind(query)
		if err := r.db.SelectContext(ctx, &addonRows, query, args...); err != nil {
			return nil, err
		}
		for _, ar := range addonRows {
			addonsByItem[ar.CartItemID] = append(addonsByItem[ar.CartItemID], model.CartItemAddonDetail{
				AddonID:   ar.AddonID,
				AddonName: ar.AddonName,
				Price:     ar.Price,
			})
		}
	}

	detail := &model.CartDetail{Cart: cart, Items: make([]model.CartItemDetail, 0, len(rows))}
	for _, row := range rows {
		item := model.CartItemDetail{
			CartItem: model.CartItem{
				ID:        row.ID,
				CartID:    row.CartID,
				ProductID: row.ProductID,
				Quantity:  row.Quantity,
			},
			ProductName:  row.ProductName,
			ProductImage: row.ProductImage,
			Addons:       addonsByItem[row.ID],
		}
		if row.VariantID.Valid {
			v := row.VariantID.Int64
			item.CartItem.VariantID = &v
		}
		if row.VariantName.Valid {
			s := row.VariantName.String
			item.VariantName = &s
		}
		if row.VariantUnit.Valid {
			s := row.VariantUnit.String
			item.VariantUnit = &s
		}
		if row.VariantPrice.Valid {
			p := row.VariantPrice.Int64
			item.VariantPrice = &p
		}
		if item.Addons == nil {
			item.Addons = []model.CartItemAddonDetail{}
		}
		detail.Items = append(detail.Items, item)
	}

	return detail, nil
}

func (r *CartRepo) AddItem(ctx context.Context, cartID string, productID int64, variantID *int64, quantity int, addonIDs []int64) (*model.CartItem, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var item model.CartItem
	err = tx.QueryRowContext(ctx,
		`INSERT INTO cart_items (cart_id, product_id, variant_id, quantity)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, cart_id, product_id, variant_id, quantity, created_at, updated_at`,
		cartID, productID, variantID, quantity,
	).Scan(&item.ID, &item.CartID, &item.ProductID, &item.VariantID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	for _, addonID := range addonIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO cart_item_addons (cart_item_id, addon_id) VALUES ($1, $2)`,
			item.ID, addonID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepo) UpdateItem(ctx context.Context, cartID string, itemID int64, quantity int) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE cart_items SET quantity=$1 WHERE id=$2 AND cart_id=$3`,
		quantity, itemID, cartID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, cartID string, itemID int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM cart_items WHERE id=$1 AND cart_id=$2`, itemID, cartID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *CartRepo) ClearCart(ctx context.Context, cartID string) error {
	// Verify cart exists — DELETE on cart_items won't fail if cart is missing
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM carts WHERE id=$1)`, cartID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return apperror.ErrNotFound
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id=$1`, cartID)
	return err
}
