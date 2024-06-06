package cart

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB создает in-memory базу данных для тестирования

func TestService_GetCartItems(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	service := NewService(repo)

	ctx := context.Background()
	userID := "test-user"

	// Добавляем тестовые данные
	_, err := db.ExecContext(ctx, `INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)`, userID, "123", 1)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	items, err := service.GetCartItems(ctx, userID)
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

func TestService_AddToCart(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	service := NewService(repo)

	ctx := context.Background()
	userID := "test-user"
	dto := AddToCartDTO{
		ProductID: "123",
		Quantity:  1,
	}

	err := service.AddToCart(ctx, userID, dto)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int
	err = db.GetContext(ctx, &count, `SELECT COUNT(*) FROM cart_items WHERE user_id=$1 AND product_id=$2`, userID, dto.ProductID)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}
