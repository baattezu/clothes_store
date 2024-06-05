package cart

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

func (r *Repository) GetCartItems(ctx context.Context, userID string) ([]CartItem, error) {
	var items []CartItem
	query := "SELECT product_id, quantity FROM cart_items WHERE user_id=$1"
	err := r.db.SelectContext(ctx, &items, query, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) AddToCart(ctx context.Context, userID string, item CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (user_id, product_id) 
		DO UPDATE SET quantity = EXCLUDED.quantity
	`
	_, err := r.db.ExecContext(ctx, query, userID, item.ProductID, item.Quantity)
	return err
}
