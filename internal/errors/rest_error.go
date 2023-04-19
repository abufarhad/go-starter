package errors

import (
	"errors"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"net/http"
)

var (
	ErrInvalidJwtSigningMethod = NewError("invalid signing method while parsing jwt")
	ErrSomethingWentWrong      = "Something went wrong"
)

type RestErr struct {
	Message     string `json:"message"`
	ReferenceNo string `json:"reference_no"`
	Detail      string `json:"detail"`
	Status      int    `json:"status"`
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

func NewBadRequestError(message string, err error) *RestErr {
	restErr := &RestErr{
		Message: message,
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

func NewNotFoundError(message string, err error) *RestErr {
	restErr := &RestErr{
		Message: message,
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

func NewAlreadyExistError(message string, err error) *RestErr {
	restErr := &RestErr{
		Message: message,
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

func NewUnauthorizedError(message string, err error) *RestErr {
	restErr := &RestErr{
		Message: message,
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
