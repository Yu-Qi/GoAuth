package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Yu-Qi/GoAuth/api"
	"github.com/Yu-Qi/GoAuth/api/middleware"
	"github.com/Yu-Qi/GoAuth/pkg/service/crypto"
	"github.com/Yu-Qi/GoAuth/pkg/service/email"
)

func main() {
	r := gin.New()
	r.Use(
		middleware.HandlePanic,
	)

	initService()

	registerAccountAPI(r)
	registerProductAPI(r)

	startServer(r)
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

func initService() {
	verificationCodeExpireSec := 600
	crypto.InitService("your-strong-password", "your-salt-string", 4096, verificationCodeExpireSec)
	email.InitService(email.NewPrintEmailService())
}

func startServer(r *gin.Engine) {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		panic("APP_PORT is not set")
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", appPort),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("listen: %s\n", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown: %v", err)
	}
}
