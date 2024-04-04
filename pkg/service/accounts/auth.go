package accounts

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Yu-Qi/GoAuth/domain"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/db"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

// RegisterParams is the parameters for registering a new account
type RegisterParams struct {
	Email    string
	Password string
}

// Register registers a new account
func Register(ctx context.Context, account *RegisterParams, verificationSvc domain.VerificationCodeService, sendEmailSvc domain.SendEmailService) (customErr *code.CustomError) {
	uid := util.UUID()
	logrus.WithFields(logrus.Fields{
		"uid":   uid,
		"email": account.Email,
	}).Debug("Registering new account")

	hashedPassword, err := util.GenerateBcryptPassword(account.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Register, GenerateBcryptPassword")
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

	verificationCode, err := verificationSvc.GenerateCode(uid)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":   uid,
			"email": account.Email,
			"error": err.Error(),
		}).Error("Register, already registered, but failed to generate verification code")
		return code.NewCustomError(code.CryptoError, http.StatusInternalServerError, err)
	}

	err = sendEmailSvc.SendEmail(account.Email, "Verification Code", verificationCode)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":   uid,
			"email": account.Email,
			"error": err.Error(),
		}).Error("Register, already registered, but failed to send email")
		return code.NewCustomError(code.SendEmailError, http.StatusInternalServerError, err)
	}
	// update sent_at
	db.UpdateAccount(ctx, uid, &domain.UpdateAccountParams{
		SentAt: util.Ptr(time.Now()),
	})

	return nil
}

// VerifyEmail verifies the email
func VerifyEmail(ctx context.Context, verificationCode string, verificationSvc domain.VerificationCodeService) (customErr *code.CustomError) {
	uid, err := verificationSvc.VerifyCode(verificationCode)
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
