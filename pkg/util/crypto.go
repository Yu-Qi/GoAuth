package util

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// default cost is 10, min is 4, max is 31
	bcryptCost = 15
)

// GenerateBcryptPassword generate bcrypt password
func GenerateBcryptPassword(password string) ([]byte, error) {
	cost := bcryptCost
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

// CompareBcryptPassword compare bcrypt password
func CompareBcryptPassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
