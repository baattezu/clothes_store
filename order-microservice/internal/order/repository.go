package order

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateOrder(ctx context.Context, order Order) error {
	tx := r.db.MustBegin()

	query := "INSERT INTO orders (user_id) VALUES ($1) RETURNING id"
	var orderID int
	err := tx.QueryRowContext(ctx, query, order.UserID).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	order.ID = orderID // Присваиваем ID объекту Order

	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity) 
			VALUES ($1, $2, $3)
		`
		_, err := tx.ExecContext(ctx, itemQuery, orderID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetOrder(ctx context.Context, orderID int) (Order, error) {
	var order Order
	query := "SELECT id, user_id, created_at FROM orders WHERE id=$1"
	err := r.db.GetContext(ctx, &order, query, orderID)
	if err != nil {
		return order, err
	}

	itemQuery := "SELECT product_id, quantity FROM order_items WHERE order_id=$1"
	err = r.db.SelectContext(ctx, &order.Items, itemQuery, orderID)
	return order, err
}
