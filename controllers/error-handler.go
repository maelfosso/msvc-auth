package controllers

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Error      error
	Message    string
	StatusCode int
}

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Message, err.StatusCode)
	}
}

func InternalError(err error, where string) *AppError {
	return &AppError{
		err,
		"An internal error occurred",
		http.StatusInternalServerError,
	}
}

func DecodingRequestBodyError(err error) *AppError {
	return &AppError{
		err,
		"Impossible to decode request data",
		http.StatusUnprocessableEntity,
	}
}

func EncodingResponseError(err error, data interface{}) *AppError {
	return &AppError{
		err,
		"impossible to encode the response",
		http.StatusInternalServerError,
	}
}

func DatabaseError(err error, message string) *AppError {
	return &AppError{
		err,
		fmt.Sprintf("Database : %s", message),
		http.StatusInternalServerError,
	}
}

func BadRequestError(err error, message string) *AppError {
	return &AppError{
		err,
		message,
		http.StatusBadRequest,
	}
}
