package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yu-Qi/GoAuth/pkg/code"
	"github.com/Yu-Qi/GoAuth/pkg/jwt"
	"github.com/Yu-Qi/GoAuth/pkg/service/accounts"
	"github.com/Yu-Qi/GoAuth/pkg/service/crypto"
	"github.com/Yu-Qi/GoAuth/pkg/service/email"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

type registerParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *registerParams) AfterValidate() error {
	if valid := util.ValidatePassword(r.Password); !valid {
		return fmt.Errorf("invalid password")
	}

	return nil
}

// Register registers a new account
func Register(c *gin.Context) {
	params := registerParams{}
	if customErr := util.ToGinContextExt(c).BindJson(&params); customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}

	customErr := accounts.Register(c, &accounts.RegisterParams{
		Email:    params.Email,
		Password: params.Password,
	}, crypto.GetService(), email.GetService())
	if customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
	})
}

type loginParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResp struct {
	AccessToken string `json:"access_token"`
}

// Login logs in an account with email and password, and returns an access token
func Login(c *gin.Context) {
	params := loginParams{}
	if customErr := util.ToGinContextExt(c).BindJson(&params); customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}

	uid, customErr := accounts.Login(c, &accounts.LoginParams{
		Email:    params.Email,
		Password: params.Password,
	})
	if customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}

	strat := jwt.NewJwtService()
	accessToken, err := strat.CreateToken(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"code":    code.CryptoError,
			"message": err,
		})
		return
	}

	resp := loginResp{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": resp,
	})
}

type verifyEmailParams struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}

// VerifyEmail verifies an email
func VerifyEmail(c *gin.Context) {
	params := verifyEmailParams{}
	if customErr := util.ToGinContextExt(c).BindJson(&params); customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}

	customErr := accounts.VerifyEmail(c, params.VerificationCode, crypto.GetService())
	if customErr != nil {
		c.JSON(customErr.HttpStatus, map[string]interface{}{
			"status":  customErr.HttpStatus,
			"code":    customErr.Code,
			"message": customErr.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
	})
}
