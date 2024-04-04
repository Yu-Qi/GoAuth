package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yu-Qi/GoAuth/pkg/code"
)

// AfterBinding is a function that gets called after the successful binding provided by GinContextExt.
// An object implementing this interface would
// gets called after the successful binding provided by GinContextExt.
type AfterBinding interface {
	AfterBinding(*gin.Context) error
}

// AfterValidate is a function that gets called after the successful binding provided by GinContextExt.
// An object implementing this interface would
// gets called after the successful binding provided by GinContextExt.
type AfterValidate interface {
	AfterValidate() error
}

// ToGinContextExt converts a gin.Context to GinContextExt.
func ToGinContextExt(ginCtx *gin.Context) *GinContextExt {
	return (*GinContextExt)(ginCtx)
}

// GinContextExt .
// Provides convenient methods for post-processing of any object.
type GinContextExt gin.Context

// Bind binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) Bind(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	if err := ctx.ShouldBind(v); err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindJson binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindJson(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	if err := ctx.ShouldBindJSON(v); err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindQuery binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindQuery(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	if err := ctx.ShouldBindQuery(v); err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindUri binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindUri(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	if err := ctx.ShouldBindUri(v); err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindHeader binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindHeader(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	if err := ctx.ShouldBindHeader(v); err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindRawBody binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindRawBody(v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

// BindFormValue binds the request body to a struct and panics if there is any error.
func (self *GinContextExt) BindFormValue(key string, v interface{}) *code.CustomError {
	ctx := (*gin.Context)(self)
	val := ctx.Request.FormValue(key)
	err := json.Unmarshal([]byte(val), v)
	if err != nil {
		return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
	}
	return self.afterBindingAndValidate(ctx, v)
}

func (self *GinContextExt) afterBindingAndValidate(ginCtx *gin.Context, v interface{}) *code.CustomError {
	if afterBinding, ok := v.(AfterBinding); ok {
		err := afterBinding.AfterBinding(ginCtx)
		if err != nil {
			return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
		}
	}

	if afterValidate, ok := v.(AfterValidate); ok {
		err := afterValidate.AfterValidate()
		if err != nil {
			return code.NewCustomError(code.ParamIncorrect, http.StatusBadRequest, err)
		}
	}

	return nil
}
