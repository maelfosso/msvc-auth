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
