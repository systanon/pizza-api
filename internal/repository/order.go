package repository

import (
	"context"
	"database/sql"

	"pizza-backend/internal/apperror"
	"pizza-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, cartID, email string, items []model.OrderItemInput, total int64) (*model.Order, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var order model.Order
	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (cart_id, email, total)
		 VALUES ($1, $2, $3)
		 RETURNING id, cart_id, email, status, total, created_at, updated_at`,
		cartID, email, total,
	).Scan(&order.ID, &order.CartID, &order.Email, &order.Status, &order.Total, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var itemID int64
		err = tx.QueryRowContext(ctx,
			`INSERT INTO order_items
			   (order_id, product_id, product_name, variant_id, variant_name, variant_price, quantity, item_total)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 RETURNING id`,
			order.ID, item.ProductID, item.ProductName,
			item.VariantID, item.VariantName, item.VariantPrice,
			item.Quantity, item.ItemTotal,
		).Scan(&itemID)
		if err != nil {
			return nil, err
		}

		for _, addon := range item.Addons {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO order_item_addons (order_item_id, addon_id, addon_name, price)
				 VALUES ($1, $2, $3, $4)`,
				itemID, addon.AddonID, addon.AddonName, addon.Price,
			); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, orderID int64) (*model.OrderDetail, error) {
	var order model.Order
	err := r.db.GetContext(ctx, &order,
		`SELECT id, cart_id, email, status, total, created_at, updated_at
		 FROM orders WHERE id = $1`, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	var items []model.OrderItem
	if err := r.db.SelectContext(ctx, &items,
		`SELECT id, order_id, product_id, product_name, variant_id,
		        variant_name, variant_price, quantity, item_total
		 FROM order_items WHERE order_id = $1 ORDER BY id`, orderID); err != nil {
		return nil, err
	}

	itemIDs := make([]int64, 0, len(items))
	for _, it := range items {
		itemIDs = append(itemIDs, it.ID)
	}

	addonsByItem := map[int64][]model.OrderItemAddon{}
	if len(itemIDs) > 0 {
		query, args, err := sqlx.In(
			`SELECT order_item_id, addon_id, addon_name, price
			 FROM order_item_addons WHERE order_item_id IN (?)`, itemIDs)
		if err != nil {
			return nil, err
		}
		var addons []model.OrderItemAddon
		if err := r.db.SelectContext(ctx, &addons, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
		for _, a := range addons {
			addonsByItem[a.OrderItemID] = append(addonsByItem[a.OrderItemID], a)
		}
	}

	detail := &model.OrderDetail{Order: order, Items: make([]model.OrderItemDetail, 0, len(items))}
	for _, it := range items {
		addons := addonsByItem[it.ID]
		if addons == nil {
			addons = []model.OrderItemAddon{}
		}
		detail.Items = append(detail.Items, model.OrderItemDetail{OrderItem: it, Addons: addons})
	}

	return detail, nil
}
