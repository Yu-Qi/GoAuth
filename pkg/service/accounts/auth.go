package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db"
	"github.com/Yu-Qi/GoAuth/pkg/log"
	"github.com/Yu-Qi/GoAuth/pkg/service/crypto"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

// RegisterParams is the parameters for registering a new account
type RegisterParams struct {
	Email    string
	Password string
}

// Register registers a new account
func Register(ctx context.Context, account *RegisterParams) (customErr *code.CustomError) {
	uid := util.UUID()
	log.DebugWithDataCtx(ctx, "Register", map[string]interface{}{
		"uid":   uid,
		"email": account.Email,
	})

	// validate email
	if !util.ValidateEmail(account.Email) {
		return
	}
	// validate password
	if !util.ValidatePassword(account.Password) {
		return
	}
	hashedPassword, err := util.GenerateBcryptPassword(account.Password)
	if err != nil {
		// TODO: 處理
		return
	}

	// check if account already exists
	if customErr := db.CreateAccount(ctx, &db.CreateAccountParams{
		UID:            uid,
		Email:          account.Email,
		HashedPassword: string(hashedPassword),
	}); customErr != nil {
		return customErr
	}

	verificationCode, err := crypto.GetService().GenerateCode(uid)
	if err != nil {
		return code.NewCustomError(code.CryptoError, http.StatusInternalServerError, err)
	}

	// TODO: send email
	fmt.Println("~~~verificationCode", verificationCode)
	return nil
}

// VerifyEmail verifies the email
func VerifyEmail(ctx context.Context, verificationCode string) (customErr *code.CustomError) {
	uid, err := crypto.GetService().VerifyCode(verificationCode)
	if err != nil {
		return code.NewCustomError(code.CryptoError, http.StatusBadRequest, err)
	}

	// activate the account
	customError := db.ActiveAccount(ctx, uid)
	if customError != nil {
		return customError
	}

	return nil
}

// LoginParams is the parameters for login
type LoginParams struct {
	Email    string
	Password string
}

// Login login an active account
func Login(ctx context.Context, params *LoginParams) (string, *code.CustomError) {
	// check if account already exists
	account, customErr := db.Login(ctx, &domain.Account{
		Email: params.Email,
	})
	if customErr != nil {
		return "", customErr
	}
	// check if password is correct
	if util.CompareBcryptPassword(account.HashedPassword, params.Password) != nil {
		return "", code.NewCustomError(code.AccountOrPasswordIncorrect, http.StatusBadRequest, fmt.Errorf("account or password incorrect"))
	}

	// check if account is active
	if !account.IsActive {
		return "", code.NewCustomError(code.AccountNotActive, http.StatusBadRequest, fmt.Errorf("account not active"))
	}

	return account.UID, nil
}
