package api

import (
	"github.com/Yu-Qi/GoAuth/api/response"
	"github.com/Yu-Qi/GoAuth/pkg/service/products"
	"github.com/gin-gonic/gin"
)

// GetRecommendations returns a list of recommended products
func GetRecommendations(ctx *gin.Context) {
	products, customErr := products.GetRecommendations(ctx)

	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	response.OK(ctx, products)
}
