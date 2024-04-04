package util

import (
	"regexp"
)

var (
	uppercasePattern   = regexp.MustCompile(`[A-Z]`)
	lowercasePattern   = regexp.MustCompile(`[a-z]`)
	specialCharPattern = regexp.MustCompile(`[\(\)\[\]\{\}<>+\-*/?,.:;"'_\\|~` + "`" + `!@#$%^&=]`)
)

// ValidatePassword check if it matches the requirement
func ValidatePassword(password string) bool {
	// rules:
	// 1. password must be at least 6 characters and no more than 16 characters
	// 2. password must contain at least one uppercase letter and one lowercase letter
	// 3. password must contain at least one special character ()[]{}<>+-*/?,.:;"'_\|~`!@#$%^&=

	if len(password) < 6 || len(password) > 16 {
		return false
	}

	return uppercasePattern.MatchString(password) &&
		lowercasePattern.MatchString(password) &&
		specialCharPattern.MatchString(password)
}
