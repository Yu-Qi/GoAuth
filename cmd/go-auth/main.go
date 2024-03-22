package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Yu-Qi/GoAuth/api"
	"github.com/Yu-Qi/GoAuth/api/middleware"
	"github.com/Yu-Qi/GoAuth/pkg/service/crypto"
)

func main() {
	r := gin.New()
	r.Use(
		middleware.HandlePanic,
	)
	crypto.Init("your-strong-password", "your-salt-string", 4096, 600)

	registerAccountAPI(r)
	registerProductAPI(r)
	appPort := os.Getenv("APP_PORT")
	_ = r.Run(":" + appPort)

}

func registerAccountAPI(r *gin.Engine) {
	r.POST("/register", api.Register)
	r.POST("/login", api.Login)
	r.POST("/verify-email", api.VerifyEmail)
}

func registerProductAPI(r *gin.Engine) {
	product := r.Group("/products", middleware.AuthToken)
	product.GET("/recommendation", api.GetRecommendations)
}
