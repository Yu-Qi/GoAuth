package code

import (
	"errors"
)

// constants
const (
	OK = 0
	// general errors
	ParamIncorrect   = 1000
	NotFound         = 1001
	DBError          = 1002
	TokenInValid     = 1003
	CacheError       = 1004
	UserNotFound     = 1005
	InvalidSignature = 1006
	JsonMarshalError = 1007
	JsonUnmarshalErr = 1008
	// business errors
	AccountAlreadyExists       = 2000
	AccountOrPasswordIncorrect = 2001
	AccountNotActive           = 2002
	AccountAlreadyActive       = 2003
	// internal errors
	CryptoError          = 3000
	InternalUnknownError = 3999
)

// define errors
var (
	ErrParamIncorrect = errors.New("param incorrect")
	ErrInternalError  = errors.New("internal error")
)

// CustomError a custom error
type CustomError struct {
	Code       int   `json:"code"`
	Error      error `json:"error"`
	HttpStatus int   `json:"http_status"`
}

// NewCustomError create a new CustomError
func NewCustomError(code, httpStatus int, err error) *CustomError {
	return &CustomError{
		Code:       code,
		Error:      err,
		HttpStatus: httpStatus,
	}
}
