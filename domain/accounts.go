package domain

// Account is a struct that represents a user account
type Account struct {
	UID            string `json:"uid"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	IsActive       bool   `json:"-"`
}
