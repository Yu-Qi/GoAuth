package products

import (
	"context"
	"net/http"
	"time"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/cache"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db"
)

const (
	ProductRecommendationCacheTTLSec = 60 * 10
)

// GetRecommendations gets product recommendations from cache or db
func GetRecommendations(ctx context.Context) ([]domain.Product, *code.CustomError) {
	// if hit cache
	if cache.Exists(ctx, cache.CacheKeyProductRecommendation) {
		v, err := cache.Get(ctx, cache.CacheKeyProductRecommendation)
		if err != nil {
			return nil, code.NewCustomError(code.CacheError, http.StatusInternalServerError, err)
		}
		return v.([]domain.Product), nil
	}

	// if not hit cache, get from db and set cache
	products, customErr := db.GetProductRecommendations(ctx)
	if customErr != nil {
		return nil, customErr
	}
	cache.Set(ctx, cache.CacheKeyProductRecommendation, products, ProductRecommendationCacheTTLSec*time.Second)

	return products, nil
}
