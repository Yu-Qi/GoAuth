package db

import (
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	// ErrDuplicateEntryCode is the error code for duplicate entry
	ErrDuplicateEntryCode = 1062
)

// IsDuplicateEntryError checks if the error is a duplicate entry error
func IsDuplicateEntryError(err error) bool {
	return mysqlErrCode(err) == ErrDuplicateEntryCode
}

// IsRecordNotFoundError checks if the error is a record not found error
func IsRecordNotFoundError(err error) bool {
	return err != nil && err == gorm.ErrRecordNotFound
}

func mysqlErrCode(err error) int {
	if err == nil {
		return 0
	}

	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return 0
	}
	return int(mysqlErr.Number)
}
