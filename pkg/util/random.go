package util

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// UUID generates a new UUID in version 4
func UUID() string {
	return uuid.NewString()
}

// RandString returns a random string of length n
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandEmail returns a random email
func RandEmail() string {
	return strings.ToLower(RandString(10)) + "@kryptogo.com"
}
