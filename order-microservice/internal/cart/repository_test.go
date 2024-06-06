package cart

import (
	"context"
	"testing"

	_ "github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	schema := `
	CREATE TABLE cart_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		product_id TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id, product_id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func TestRepository_GetCartItems(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	ctx := context.Background()
	userID := "test-user"

	_, err := db.ExecContext(ctx, `INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)`, userID, "123", 1)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	items, err := repo.GetCartItems(ctx, userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if items[0].ProductID != "123" || items[0].Quantity != 1 {
		t.Errorf("Unexpected item values: %+v", items[0])
	}
}

func TestRepository_AddToCart(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	ctx := context.Background()
	userID := "test-user"
	item := CartItem{
		ProductID: "123",
		Quantity:  1,
	}

	err := repo.AddToCart(ctx, userID, item)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int
	err = db.GetContext(ctx, &count, `SELECT COUNT(*) FROM cart_items WHERE user_id=$1 AND product_id=$2`, userID, item.ProductID)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}
