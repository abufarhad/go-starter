package errors

import (
	"errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"net/http"
)

var (
	ErrInvalidEmail            = NewError("invalid email")
	ErrInvalidPassword         = NewError("invalid password")
	ErrInvalidEmailOrPassword  = NewError("invalid email or password")
	ErrCreateJwt               = NewError("failed to create JWT token")
	ErrParseJwt                = NewError("failed to parse JWT token")
	ErrInvalidRefreshToken     = NewError("invalid refresh_token")
	ErrInvalidRefreshUuid      = NewError("invalid refresh_uuid")
	ErrInvalidAccessToken      = NewError("invalid access_token")
	ErrAccessTokenSign         = NewError("failed to sign access_token")
	ErrRefreshTokenSign        = NewError("failed to sign refresh_token")
	ErrNoContextUser           = NewError("failed to get auth from context")
	ErrInvalidUserType         = NewError("invalid user type")
	ErrInvalidJwtSigningMethod = NewError("invalid signing method while parsing jwt")
	ErrSomethingWentWrong      = "Something went wrong"
	ErrRecordNotFound          = "Record not found"
	ErrBadRequest              = "Bad request"
	ErrAlreadyExist            = "Record already exist"
	ErrUnauthorizedError       = "Unauthorized error"
	ErrRecordNotValid          = "invalid parameters, check email or password"
)

type RestErr struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Status  int    `json:"status"`
}

func (err *RestErr) Error() string {
	return err.Message
}

func NewError(msg string) error {
	return errors.New(msg)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}

func NewInternalServerError(err error) *RestErr {
	restErr := &RestErr{
		Message: ErrSomethingWentWrong,
		Status:  http.StatusInternalServerError,
	}

	if err != nil {
		logger.ErrorAsJson(err.Error(), err)
		restErr.Detail = err.Error()
	} else {
		logger.ErrorAsJson(restErr.Error(), restErr)
	}

	return restErr
}

func NewBadRequestError(err error) *RestErr {
	restErr := &RestErr{
		Message: ErrBadRequest,
		Status:  http.StatusBadRequest,
	}

	if err != nil {
		logger.ErrorAsJson(err.Error(), err)
		restErr.Detail = err.Error()
	} else {
		logger.ErrorAsJson(restErr.Error(), restErr)
	}

	return restErr
}

func NewNotFoundError(err error) *RestErr {
	restErr := &RestErr{
		Message: ErrRecordNotFound,
		Status:  http.StatusNotFound,
	}

	if err != nil {
		logger.ErrorAsJson(err.Error(), err)
		restErr.Detail = err.Error()
	} else {
		logger.ErrorAsJson(restErr.Error(), restErr)
	}

	return restErr
}

func NewAlreadyExistError(err error) *RestErr {
	restErr := &RestErr{
		Message: ErrAlreadyExist,
		Status:  http.StatusConflict,
	}

	if err != nil {
		logger.ErrorAsJson(err.Error(), err)
		restErr.Detail = err.Error()
	} else {
		logger.ErrorAsJson(restErr.Error(), restErr)
	}

	return restErr
}

func NewUnauthorizedError(err error) *RestErr {
	restErr := &RestErr{
		Message: ErrUnauthorizedError,
		Status:  http.StatusUnauthorized,
	}

	if err != nil {
		logger.ErrorAsJson(err.Error(), err)
		restErr.Detail = err.Error()
	} else {
		logger.ErrorAsJson(restErr.Error(), restErr)
	}

	return restErr
}
