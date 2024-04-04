package api

import (
	"net/http"

	"github.com/Yu-Qi/GoAuth/pkg/service/products"
	"github.com/gin-gonic/gin"
)

// GetRecommendations returns a list of recommended products
func GetRecommendations(c *gin.Context) {
	products, customErr := products.GetRecommendations(c, c.GetString("uid"))
	if customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": products,
	})
}
