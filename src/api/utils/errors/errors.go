package errors

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type ApiError interface {
	Status() int
	Message() string
	Error()   string
}

type apiError struct {
	status  int    `json:"status"`
	message string `json:"message"`
	error   string `json:"error"`
}

func (e *apiError) Status() int {
	return e.status
}

func (e *apiError) Message() string {
	return e.message
}

func (e *apiError) Error() string {
	return e.error
}

func NewApiError(statusCode int, message string) ApiError{
	return &apiError{
		status:  statusCode,
		message: message,

	}
}
func NewApiErrorFromBytes(b []byte) (ApiError, error){
	var result apiError
	log.Printf("errors.go::bytes: %c", b)
	if err := json.Unmarshal(b, &result); err != nil {
		log.Printf("errors.go::NewApiErrorFromBytes:Status: %d", result.Status())
		log.Printf("errors.go::NewApiErrorFromBytes:Message: %s ", result.Message())

		return nil, errors.New("invalid json body")
	}

	return &result, nil
}
func NewNotFoundError(message string) ApiError {
	return &apiError{
		status:  http.StatusNotFound,
		message: message,

	}
}

func NewInternalServerError(message string) ApiError {
	return &apiError{
		status:  http.StatusInternalServerError,
		message: message,

	}
}

func NewBadRequestError(message string) ApiError {
	return &apiError{
		status:  http.StatusBadRequest,
		message: message,

	}
}
