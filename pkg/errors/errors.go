package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternal          = errors.New("internal error")
	ErrTimeout           = errors.New("operation timeout")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrValidation        = errors.New("validation failed")
)

type AppError struct {
	Err     error
	Message string
	Code    string
	Details map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

func (e *AppError) WithCode(code string) *AppError {
	e.Code = code
	return e
}

func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

func ValidationError(message string) *AppError {
	return New(ErrValidation, message)
}

func DatabaseError(message string, err error) *AppError {
	return New(err, message).WithCode("DB_ERROR")
}

func NotFoundError(message string) *AppError {
	return New(ErrNotFound, message).WithCode("NOT_FOUND")
}

func InternalError(message string, err error) *AppError {
	return New(err, message).WithCode("INTERNAL_ERROR")
}
