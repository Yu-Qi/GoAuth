package config

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GetString .
func GetString(key string) string {
	return os.Getenv(key)
}

// GetInt .
func GetInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(err)
	}
	return val
}
