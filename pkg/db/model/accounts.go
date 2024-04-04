package model

import "time"

// TableNameAccount is the table name of <accounts>
const TableNameAccount = "accounts"

// Account mapped from table <accounts>
type Account struct {
	UID            string     `gorm:"column:uid;type:varchar(36);not null;primaryKey"`
	Email          string     `gorm:"column:email;type:varchar(256);not null;uniqueIndex:idx_accounts"`
	HashedPassword string     `gorm:"column:hashed_password;type:varchar(72);not null"`
	IsActive       bool       `gorm:"column:is_active;type:tinyint(1);not null;default:0"`
	SentAt         *time.Time `gorm:"column:sent_at;type:timestamp;"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeleteAt       *time.Time `gorm:"column:delete_at;type:timestamp"`
}

// TableName Account's table name
func (*Account) TableName() string {
	return TableNameAccount
}
