package xerror

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	ErrCodeNotFound            = "not_found"
	ErrCodeBadRequest          = "bad_request"
	ErrCodeInternalServerError = "internal_server_error"
	ErrCodeForbidden           = "forbidden"
	ErrCodeUnauthorized        = "unauthorized"
)

var _ error = (*Error)(nil)

// Error is an API error, an error that has a special handling during the responses.
type Error struct {
	errorCode string
	args      map[string]any
}

func (e Error) ErrorCode() string {
	return e.errorCode
}

func (e Error) ErrorArgs() map[string]any {
	return e.args
}

func (e Error) Error() string {
	var args []string
	for key, value := range e.args {
		args = append(args, fmt.Sprintf("%s=%v", key, value))
	}
	slices.Sort(args)
	return fmt.Sprintf("%s: %s", e.errorCode, strings.Join(args, ", "))
}

type ErrorParam = func(e *Error)

func ErrParam[T any](key string, value T) ErrorParam {
	return func(e *Error) {
		e.args[key] = value
	}
}

// MakeError creates an error.
func MakeError(code string, params ...ErrorParam) error {
	e := Error{
		errorCode: code,
		args:      make(map[string]any),
	}
	for _, p := range params {
		p(&e)
	}
	return e
}

// MakeBadRequestError creates a bad request error.
func MakeBadRequestError(params ...ErrorParam) error {
	e := Error{
		errorCode: ErrCodeBadRequest,
		args:      make(map[string]any),
	}
	for _, p := range params {
		p(&e)
	}
	return e
}

// MakeNotFoundError creates a not found error.
func MakeNotFoundError(params ...ErrorParam) error {
	e := Error{
		errorCode: ErrCodeNotFound,
		args:      make(map[string]any),
	}
	for _, p := range params {
		p(&e)
	}
	return e
}

// MakeForbiddenError creates a forbidden error.
func MakeForbiddenError(params ...ErrorParam) error {
	e := Error{
		errorCode: ErrCodeForbidden,
		args:      make(map[string]any),
	}
	for _, p := range params {
		p(&e)
	}
	return e
}

// MakeUnauthorizedError creates a unauthorized error.
func MakeUnauthorizedError(params ...ErrorParam) error {
	e := Error{
		errorCode: ErrCodeUnauthorized,
		args:      make(map[string]any),
	}
	for _, p := range params {
		p(&e)
	}
	return e
}

// IsErrNotFound check if the given error is a not found error.
func IsErrNotFound(err error) bool {
	return IsErrCode(err, ErrCodeNotFound)
}

// IsErrCode check if the given error has the given code.
func IsErrCode(err error, code string) bool {
	var apiError Error
	if errors.As(err, &apiError) {
		return apiError.ErrorCode() == code
	}
	return false
}
