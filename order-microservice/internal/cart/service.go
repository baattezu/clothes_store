package cart

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

func (s *Service) GetCartItems(ctx context.Context, userID string) ([]CartItem, error) {
	return s.repo.GetCartItems(ctx, userID)
}

func (s *Service) AddToCart(ctx context.Context, userID string, dto AddToCartDTO) error {
	item := CartItem{
		ProductID: dto.ProductID,
		Quantity:  dto.Quantity,
	}
	return s.repo.AddToCart(ctx, userID, item)
}
