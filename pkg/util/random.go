package util

import (
	"github.com/google/uuid"
)

// UUID generates a new UUID in version 4
func UUID() string {
	return uuid.NewString()
}
