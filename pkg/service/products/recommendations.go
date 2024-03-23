package products

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/cache"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db"
)

var (
	g singleflight.Group
)

const (
	ProductRecommendationCacheTTLSec = 60 * 10
)

// GetRecommendations gets product recommendations from cache or db
func GetRecommendations(ctx context.Context) (products []domain.Product, customErr *code.CustomError) {
	// if hit cache
	if cache.Exists(ctx, cache.CacheKeyProductRecommendation) {
		v, err := cache.Get(ctx, cache.CacheKeyProductRecommendation)
		if err != nil {
			return nil, code.NewCustomError(code.CacheError, http.StatusInternalServerError, err)
		}
		err = json.Unmarshal([]byte(v.(string)), &products)
		if err != nil {
			return nil, code.NewCustomError(code.JsonUnmarshalErr, http.StatusInternalServerError, err)
		}
		return products, nil
	}

	// if not hit cache, get from db and set cache
	// use singleflight to avoid cache breakdown.
	data, err, _ := g.Do("", func() (any, error) {
		products, customErr = db.GetProductRecommendations(ctx)
		if customErr != nil {
			return nil, customErr.Error
		}

		err := cache.SetWithObject(ctx, cache.CacheKeyProductRecommendation, products, ProductRecommendationCacheTTLSec*time.Second)
		if err != nil {
			customErr = code.NewCustomError(code.CacheError, http.StatusInternalServerError, err)
			return nil, err
		}
		return products, nil
	})
	if err != nil {
		return
	}

	return data.([]domain.Product), nil
}
