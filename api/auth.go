package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yu-Qi/GoAuth/api/response"
	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/jwt"
	"github.com/Yu-Qi/GoAuth/pkg/service/accounts"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

type registerParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register registers a new account
func Register(ctx *gin.Context) {
	params := registerParams{}
	customErr := util.ToGinContextExt(ctx).BindJson(&params)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	customErr = accounts.Register(ctx, &accounts.RegisterParams{
		Email:    params.Email,
		Password: params.Password,
	})
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	response.OK(ctx, nil)
}

type loginParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	AccessToken string `json:"access_token"`
}

// Login logs in an account with email and password, and returns an access token
func Login(ctx *gin.Context) {
	params := loginParams{}
	customErr := util.ToGinContextExt(ctx).BindJson(&params)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	uid, customErr := accounts.Login(ctx, &accounts.LoginParams{
		Email:    params.Email,
		Password: params.Password,
	})
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	strat := jwt.NewJwtService()
	accessToken, err := strat.CreateToken(uid)
	if err != nil {
		response.CustomError(ctx, code.NewCustomError(code.CryptoError, http.StatusInternalServerError, err))
		return
	}

	resp := loginResp{
		AccessToken: accessToken,
	}
	response.OK(ctx, resp)
}

type verifyEmailParams struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}

// VerifyEmail verifies an email
func VerifyEmail(ctx *gin.Context) {
	params := verifyEmailParams{}
	customErr := util.ToGinContextExt(ctx).BindJson(&params)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}
	customErr = accounts.VerifyEmail(ctx, params.VerificationCode)
	if customErr != nil {
		response.CustomError(ctx, customErr)
		return
	}

	response.OK(ctx, nil)
}
