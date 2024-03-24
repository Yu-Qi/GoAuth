package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/Yu-Qi/GoAuth/api/middleware"
	"github.com/Yu-Qi/GoAuth/pkg/db"
	"github.com/Yu-Qi/GoAuth/pkg/db/model"
	jwtSvc "github.com/Yu-Qi/GoAuth/pkg/jwt"
	"github.com/Yu-Qi/GoAuth/pkg/util"
)

type getRecommendationsSuite struct {
	suite.Suite
	Url     string
	Request func(headers http.Header) (httpStatus int, responseBody []byte, err error)
}

func (suite *getRecommendationsSuite) SetupSuite() {
	suite.Url = fmt.Sprintf("/products/recommendation")
	suite.Request = func(headers http.Header) (httpStatus int, responseBody []byte, err error) {
		return util.GetWithHeaderForTest(suite.Url, headers, middleware.AuthToken, GetRecommendations)

	}
}

func TestGetRecommendations(t *testing.T) {
	suite.Run(t, new(getRecommendationsSuite))
}

func (suite *getRecommendationsSuite) TestNormal() {
	// setup a new account in the database
	uid := util.UUID()
	email := util.RandEmail()
	password := "Password1!" + util.RandString(3)
	hashedPassword, err := util.GenerateBcryptPassword(password)
	assert.Nil(suite.T(), err)
	db.Get().Create(&model.Account{UID: uid, Email: email, HashedPassword: string(hashedPassword), IsActive: true})

	// generate a token
	strat := jwtSvc.NewJwtService()
	token, err := strat.CreateToken(uid)
	assert.Nil(suite.T(), err)
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	httpStatus, respBody, err := suite.Request(headers)
	var resp struct {
		Code int `json:"code"`
		Data []struct {
			ProductID int    `json:"product_id"`
			Name      string `json:"name"`
		} `json:"data"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, resp.Code)
	assert.Equal(suite.T(), 1, len(resp.Data))
	assert.Equal(suite.T(), 1, resp.Data[0].ProductID)
	assert.Equal(suite.T(), "product1", resp.Data[0].Name)
}

func (suite *getRecommendationsSuite) TestInvalidToken() {
	headers := http.Header{
		"Authorization": []string{"Bearer invalid_token"},
	}
	httpStatus, respBody, err := suite.Request(headers)
	var resp struct {
		Code int `json:"code"`
		Data []struct {
			ProductID int    `json:"product_id"`
			Name      string `json:"name"`
		} `json:"data"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1003, resp.Code)
}

func (suite *getRecommendationsSuite) TestNoToken() {
	headers := http.Header{}
	httpStatus, respBody, err := suite.Request(headers)
	var resp struct {
		Code int `json:"code"`
		Data []struct {
			ProductID int    `json:"product_id"`
			Name      string `json:"name"`
		} `json:"data"`
	}
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, httpStatus)
	err = json.Unmarshal(respBody, &resp)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1003, resp.Code)
}
