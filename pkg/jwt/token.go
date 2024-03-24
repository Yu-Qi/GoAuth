package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TokenJwtStrategy is the JWT strategy for custom token
type TokenJwtStrategy struct {
	secretKey string
}

// Parse parses the JWT token
func (s *TokenJwtStrategy) Parse(tokenString string) (any, *code.CustomError) {
	var claimsI interface{}
	var token *jwt.Token
	var err error
	claimsI = &jwt.StandardClaims{}
	claims := claimsI.(*jwt.StandardClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return claimsI, code.NewCustomError(code.TokenInValid, http.StatusUnauthorized, err)
	}

	if !token.Valid {
		return claimsI, code.NewCustomError(code.TokenInValid, http.StatusUnauthorized, fmt.Errorf("Invalid JWT Token"))
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return claimsI, code.NewCustomError(code.TokenExpired, http.StatusUnauthorized, fmt.Errorf("Token expired"))
	}

	return claimsI, nil
}

// CreateToken creates the JWT token
func (s *TokenJwtStrategy) CreateToken(data any) (string, error) {
	now := time.Now()
	uid, ok := data.(string)
	if !ok {
		return "", fmt.Errorf("invalid data type")
	}
	expiresAt := now.Add(time.Duration(config.GetInt("ACCESS_TOKEN_EXP_MINUTES")) * time.Minute)
	tokenClaims := jwt.StandardClaims{
		Issuer:    "Alan chen",
		Subject:   uid,
		Audience:  "https://alanchen.com",
		ExpiresAt: expiresAt.Unix(),
		NotBefore: now.Unix(),
		IssuedAt:  now.Unix(),
		Id:        uuid.New().String(),
	}
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	token, err := jwtClaims.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
