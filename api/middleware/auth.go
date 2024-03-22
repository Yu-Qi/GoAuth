package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/Yu-Qi/GoAuth/api/response"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	jwtSvc "github.com/Yu-Qi/GoAuth/pkg/jwt"
)

// AuthToken is the middleware to authenticate the token
func AuthToken(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")

	parts := strings.SplitN(authorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		response.ErrorWithMsg(ctx, http.StatusUnauthorized, code.TokenInValid, "invalid auth token")
		return
	}
	token := parts[1]
	strat := jwtSvc.NewJwtService()
	claimsI, customErr := strat.Parse(token)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	claims := claimsI.(*jwt.StandardClaims)

	userID := claims.Subject
	ctx.Set("uid", userID)
	ctx.Next()
	return
}
