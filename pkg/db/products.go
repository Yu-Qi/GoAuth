package db

import (
	"context"
	"time"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/code"
)

// GetProductRecommendations gets product recommendations
func GetProductRecommendations(ctx context.Context) ([]domain.Product, *code.CustomError) {
	time.Sleep(3 * time.Second) // simulate db slow query
	return []domain.Product{
		{
			ID:   1,
			Name: "product1",
		},
	}, nil
}
