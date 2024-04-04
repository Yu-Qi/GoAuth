package domain

import "time"

// Account is a struct that represents a user account
type Account struct {
	UID            string     `json:"uid"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"-"`
	IsActive       bool       `json:"-"`
	SentAt         *time.Time `json:"-"`
}

// UpdateAccountParams is the parameters for updating an account
type UpdateAccountParams struct {
	SentAt *time.Time
}
