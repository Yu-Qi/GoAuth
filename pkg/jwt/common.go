//go:generate go-enum
package jwt

import (
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/config"
)

var (
	jwtTokenSecret string
)

func init() {
	jwtTokenSecret = config.GetString("JWT_TOKEN_SECRET")
	if jwtTokenSecret == "" {
		panic("JWT_SECRET_KEY is empty")
	}
}

// Strategy is the interface for jwt strategy
type Strategy interface {
	Parse(tokenString string) (any, *code.CustomError)
	CreateToken(data any) (string, error)
}

// NewJwtService returns a JwtStrategy
func NewJwtService() Strategy {
	return &TokenJwtStrategy{secretKey: jwtTokenSecret}
}
