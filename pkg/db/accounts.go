package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateAccountParams is the parameters for creating an account
type CreateAccountParams struct {
	UID            string
	Email          string
	HashedPassword string
}

// CreateAccount creates a new account
func CreateAccount(ctx context.Context, params *CreateAccountParams) *code.CustomError {
	err := GetWith(ctx).Create(&model.Account{
		UID:            params.UID,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
	}).Error
	if IsDuplicateEntryError(err) {
		return code.NewCustomError(code.AccountAlreadyExists, http.StatusBadRequest, err)
	} else if err != nil {
		return code.NewCustomError(code.DBError, http.StatusInternalServerError, err)
	}

	return nil
}

// Login login an active account
func Login(ctx context.Context, params *domain.Account) (*domain.Account, *code.CustomError) {
	// check if account already exists
	account := &model.Account{}
	err := GetWith(ctx).
		Where("email = ?", params.Email).
		First(account).Error

	if IsRecordNotFoundError(err) {
		return nil, code.NewCustomError(code.AccountOrPasswordIncorrect, http.StatusBadRequest, fmt.Errorf("account or password incorrect"))
	} else if err != nil {
		return nil, code.NewCustomError(code.DBError, http.StatusInternalServerError, err)
	}

	return &domain.Account{
		UID:            account.UID,
		Email:          account.Email,
		HashedPassword: account.HashedPassword,
		IsActive:       account.IsActive,
	}, nil
}

// ActiveAccount activates an account
func ActiveAccount(ctx context.Context, uid string) *code.CustomError {
	httpStatus := http.StatusInternalServerError
	errCode := code.DBError
	err := GetWith(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("uid = ?", uid).
			First(&model.Account{}).Error
		if IsRecordNotFoundError(err) {
			httpStatus = http.StatusBadRequest
			errCode = code.UserNotFound
			return err
		} else if err != nil {
			return err
		}

		query := GetWith(ctx).
			Model(&model.Account{}).
			Where("uid = ?", uid).
			Update("is_active", true)
		if query.Error != nil {
			return err
		}
		if query.RowsAffected == 0 {
			httpStatus = http.StatusBadRequest
			errCode = code.AccountAlreadyActive
			return fmt.Errorf("account already active")
		}

		return nil
	})

	if err != nil {
		return code.NewCustomError(errCode, httpStatus, err)
	}
	return nil
}

// UserExists checks if a user exists and is active
func UserExists(ctx context.Context, uid string) *code.CustomError {
	account := &model.Account{}
	err := GetWith(ctx).
		Where("uid = ?", uid).
		First(account).Error
	if err != nil {
		if IsRecordNotFoundError(err) {
			return code.NewCustomError(code.UserNotFound, http.StatusNotFound, err)
		}
		return code.NewCustomError(code.DBError, http.StatusInternalServerError, err)
	}
	if !account.IsActive {
		return code.NewCustomError(code.AccountNotActive, http.StatusBadRequest, fmt.Errorf("account not active"))
	}
	return nil
}

// UpdateAccount updates an account
func UpdateAccount(ctx context.Context, uid string, params *domain.UpdateAccountParams) *code.CustomError {
	httpStatus := http.StatusInternalServerError
	errCode := code.DBError
	err := GetWith(ctx).Transaction(func(tx *gorm.DB) error {
		account := model.Account{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("uid = ?", uid).
			First(&account).Error
		if err != nil {
			if IsRecordNotFoundError(err) {
				httpStatus = http.StatusBadRequest
				errCode = code.UserNotFound
			}
			return err
		}
		// check params
		if params.SentAt != nil {
			account.SentAt = params.SentAt
		}
		// update account
		err = tx.Updates(&account).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return code.NewCustomError(errCode, httpStatus, err)
	}
	return nil
}
