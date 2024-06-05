package order

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userID string, dto CreateOrderDTO) (int, error) {
	order := Order{
		UserID: userID,
		Items:  dto.Items,
	}
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (s *Service) GetOrder(ctx context.Context, orderID int) (Order, error) {
	return s.repo.GetOrder(ctx, orderID)
}
