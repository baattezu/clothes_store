package order

type CreateOrderDTO struct {
	UserID string      `json:"user_id"`
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID string `db:"product_id" json:"product_id"`
	Quantity  int    `db:"quantity" json:"quantity"`
}

type Order struct {
	ID        int         `db:"id" json:"id"`
	UserID    string      `db:"user_id" json:"user_id"`
	Items     []OrderItem `json:"items"`
	CreatedAt string      `db:"created_at" json:"created_at"`
}
