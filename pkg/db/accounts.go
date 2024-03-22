package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db/model"
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
	err := GetWith(ctx).
		Model(&model.Account{}).
		Where("uid = ?", uid).
		First(&model.Account{}).Error
	if IsRecordNotFoundError(err) {
		return code.NewCustomError(code.UserNotFound, http.StatusBadRequest, err)
	} else if err != nil {
		return code.NewCustomError(code.DBError, http.StatusInternalServerError, err)
	}

	query := GetWith(ctx).
		Model(&model.Account{}).
		Where("uid = ?", uid).
		Update("is_active", true)
	if query.Error != nil {
		return code.NewCustomError(code.DBError, http.StatusInternalServerError, query.Error)
	}
	if query.RowsAffected == 0 {
		return code.NewCustomError(code.AccountAlreadyActive, http.StatusBadRequest, fmt.Errorf("account already active"))
	}

	return nil
}
