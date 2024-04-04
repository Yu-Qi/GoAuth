package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Yu-Qi/GoAuth/pkg/db"
	"github.com/Yu-Qi/GoAuth/pkg/db/model"
	"github.com/Yu-Qi/GoAuth/pkg/service/crypto"
	"github.com/Yu-Qi/GoAuth/pkg/service/email"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

type registerSuite struct {
	suite.Suite
	Url     string
	Request func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error)
}

func (suite *registerSuite) SetupSuite() {
	suite.Url = fmt.Sprintf("/register")
	suite.Request = func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error) {
		return util.PostForTest(suite.Url, body, Register)
	}
	// dependency injection
	verificationCodeExpireSec := 600
	email.InitService(email.NewPrintEmailService())
	crypto.InitService("your-strong-password", "your-salt-string", 4096, verificationCodeExpireSec)
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(registerSuite))
}

func (suite *registerSuite) TestNormal() {
	body := map[string]interface{}{
		"email":    util.RandEmail(),
		"password": "Password1!",
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, resp.Code)
}

func (suite *registerSuite) TestWrongParameter() {
	tests := []struct {
		name       string
		body       map[string]interface{}
		httpStatus int
		errCode    int
	}{
		{
			name: "Missing email",
			body: map[string]interface{}{
				"password": "Password1!",
			},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name: "Missing password",
			body: map[string]interface{}{
				"email": util.RandEmail(),
			},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Missing email and password",
			body:       map[string]interface{}{},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
	}
	var resp struct {
		Code int `json:"code"`
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			httpStatus, respBody, err := suite.Request(tt.body)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.httpStatus, httpStatus)
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.errCode, resp.Code)
		})
	}
}

func (suite *registerSuite) TestDuplicateEmail() {
	body := map[string]interface{}{
		"email":    util.RandEmail(),
		"password": "Password1!",
	}
	// first request
	_, _, _ = suite.Request(body)
	// second request with the same email
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2000, resp.Code)
}

func (suite *registerSuite) TestInvalidEmail() {
	body := map[string]interface{}{
		"email":    "invalid-email",
		"password": "Password1!",
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1000, resp.Code)
}

func (suite *registerSuite) TestInvalidPassword() {
	tests := []struct {
		name       string
		password   string
		httpStatus int
		errCode    int
	}{
		{
			name:       "Password too short",
			password:   "Pass1",
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Password too long",
			password:   "Password123!Password123!",
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Password missing uppercase",
			password:   "password123!",
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Password missing lowercase",
			password:   "PASSWORD123!",
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Password missing special character",
			password:   "Password123",
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
	}
	var resp struct {
		Code int `json:"code"`
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			body := map[string]interface{}{
				"email":    util.RandEmail(),
				"password": tt.password,
			}
			httpStatus, respBody, err := suite.Request(body)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.httpStatus, httpStatus)
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.errCode, resp.Code)
		})
	}
}

type loginSuite struct {
	suite.Suite
	Url      string
	Request  func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error)
	Email    string
	Password string
}

func (suite *loginSuite) SetupSuite() {
	suite.Url = fmt.Sprintf("/login")
	suite.Request = func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error) {
		return util.PostForTest(suite.Url, body, Login)
	}

	// setup a new account in the database
	suite.Email = util.RandEmail()
	suite.Password = "Password1!" + util.RandString(3)
	hashedPassword, err := util.GenerateBcryptPassword(suite.Password)
	if err != nil {
		panic(err)
	}
	db.Get().Create(&model.Account{UID: util.UUID(), Email: suite.Email, HashedPassword: string(hashedPassword), IsActive: true})

	// dependency injection
	verificationCodeExpireSec := 600
	crypto.InitService("your-strong-password", "your-salt-string", 4096, verificationCodeExpireSec)
}

func TestLogin(t *testing.T) {
	suite.Run(t, new(loginSuite))
}

func (suite *loginSuite) TestNormal() {
	body := map[string]interface{}{
		"email":    suite.Email,
		"password": suite.Password,
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, resp.Code)
	assert.NotEmpty(suite.T(), resp.Data.AccessToken)
}

func (suite *loginSuite) TestWrongParameter() {
	tests := []struct {
		name       string
		body       map[string]interface{}
		httpStatus int
		errCode    int
	}{
		{
			name: "Missing email",
			body: map[string]interface{}{
				"password": suite.Password,
			},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name: "Missing password",
			body: map[string]interface{}{
				"email": suite.Email,
			},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
		{
			name:       "Missing email and password",
			body:       map[string]interface{}{},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
	}
	var resp struct {
		Code int `json:"code"`
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			httpStatus, respBody, err := suite.Request(tt.body)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.httpStatus, httpStatus)
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.errCode, resp.Code)
		})
	}
}

func (suite *loginSuite) TestAccountOrPasswordIncorrect() {
	tests := []struct {
		name       string
		body       map[string]interface{}
		httpStatus int
		errCode    int
	}{
		{
			name: "Wrong email",
			body: map[string]interface{}{
				"email":    "wrong-email",
				"password": suite.Password,
			},
			httpStatus: http.StatusBadRequest,
			errCode:    2001,
		},
		{
			name: "Wrong password",
			body: map[string]interface{}{
				"email":    suite.Email,
				"password": "wrong-password",
			},
			httpStatus: http.StatusBadRequest,
			errCode:    2001,
		},
	}

	var resp struct {
		Code int `json:"code"`
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			httpStatus, respBody, err := suite.Request(tt.body)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.httpStatus, httpStatus)
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.errCode, resp.Code)
		})
	}
}

func (suite *loginSuite) TestAccountNotActive() {
	// create an inactive account
	inactiveEmail := util.RandEmail()
	inactivePassword := "Password1!" + util.RandString(3)
	hashedPassword, err := util.GenerateBcryptPassword(inactivePassword)
	assert.Nil(suite.T(), err)
	db.Get().Create(&model.Account{UID: util.UUID(), Email: inactiveEmail, HashedPassword: string(hashedPassword), IsActive: false})

	body := map[string]interface{}{
		"email":    inactiveEmail,
		"password": inactivePassword,
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2002, resp.Code)
}

type verifyEmailSuite struct {
	suite.Suite
	Url      string
	Request  func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error)
	UID      string
	Email    string
	Password string
}

func (suite *verifyEmailSuite) SetupSuite() {
	suite.Url = fmt.Sprintf("/verify-email")
	suite.Request = func(body map[string]interface{}) (httpStatus int, responseBody []byte, err error) {
		return util.PostForTest(suite.Url, body, VerifyEmail)
	}

	// setup a new account in the database
	suite.UID = util.UUID()
	suite.Email = util.RandEmail()
	suite.Password = "Password1!" + util.RandString(3)
	hashedPassword, err := util.GenerateBcryptPassword(suite.Password)
	if err != nil {
		panic(err)
	}
	db.Get().Create(&model.Account{UID: suite.UID, Email: suite.Email, HashedPassword: string(hashedPassword), IsActive: true})

	// dependency injection
	verificationCodeExpireSec := 600
	crypto.InitService("your-strong-password", "your-salt-string", 4096, verificationCodeExpireSec)
}

func TestVerifyEmail(t *testing.T) {
	suite.Run(t, new(verifyEmailSuite))
}

func (suite *verifyEmailSuite) TestNormal() {
	verificationCode, err := crypto.GetService().GenerateCode(suite.UID)
	assert.Nil(suite.T(), err)
	body := map[string]interface{}{
		"verification_code": verificationCode,
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, resp.Code)
}

func (suite *verifyEmailSuite) TestWrongParameter() {
	tests := []struct {
		name       string
		body       map[string]interface{}
		httpStatus int
		errCode    int
	}{
		{
			name:       "Missing verification code",
			body:       map[string]interface{}{},
			httpStatus: http.StatusBadRequest,
			errCode:    1000,
		},
	}
	var resp struct {
		Code int `json:"code"`
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			httpStatus, respBody, err := suite.Request(tt.body)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.httpStatus, httpStatus)
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(suite.T(), err)
			assert.Equal(suite.T(), tt.errCode, resp.Code)
		})
	}
}

func (suite *verifyEmailSuite) TestCryptoError() {
	body := map[string]interface{}{
		"verification_code": "invalid-code",
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3000, resp.Code)
}

func (suite *verifyEmailSuite) TestUserNotFound() {
	verificationCode, err := crypto.GetService().GenerateCode(util.UUID())
	assert.Nil(suite.T(), err)
	body := map[string]interface{}{
		"verification_code": verificationCode,
	}
	httpStatus, respBody, err := suite.Request(body)
	var resp struct {
		Code int `json:"code"`
	}

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1005, resp.Code)
}
