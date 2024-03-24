package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/Yu-Qi/GoAuth/pkg/code"
	jwtSvc "github.com/Yu-Qi/GoAuth/pkg/jwt"
)

// AuthToken is the middleware to authenticate the token
func AuthToken(c *gin.Context) {
	authorization := c.GetHeader("Authorization")

	parts := strings.SplitN(authorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"code":    code.TokenInValid,
			"message": "invalid auth token",
		})
		return
	}
	token := parts[1]
	strat := jwtSvc.NewJwtService()
	claimsI, customErr := strat.Parse(token)
	if customErr != nil {
		c.AbortWithStatusJSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}
	claims := claimsI.(*jwt.StandardClaims)

	userID := claims.Subject
	c.Set("uid", userID)
	c.Next()
	return
}
