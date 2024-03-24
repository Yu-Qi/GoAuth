package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// PostForTest sends a POST request to the given URL with the given body. Put the route handler functions to last handleFuncs
func PostForTest(url string, body map[string]interface{}, handleFuncs ...gin.HandlerFunc) (httpStatus int, responseBody []byte, err error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST(url, handleFuncs...)
	r.ServeHTTP(w, req)

	httpStatus = w.Code
	responseBody = w.Body.Bytes()
	return
}

// GetWithHeaderForTest sends a GET request to the given URL with the given header. Put the route handler functions to last handleFuncs
func GetWithHeaderForTest(url string, headers http.Header, handleFuncs ...gin.HandlerFunc) (httpStatus int, responseBody []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header = headers

	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET(url, handleFuncs...)

	r.ServeHTTP(w, req)

	httpStatus = w.Code
	responseBody = w.Body.Bytes()
	return
}
