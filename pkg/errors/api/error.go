package api

import (
	"log/slog"
	"net/http"
)

type Error interface {
	error
	StatusCode() int
}

type HandlerError struct {
	Method string
	URL    string
	Err    error
}

func (e *HandlerError) Error() string {
	return e.Err.Error()
}

func (e *HandlerError) Unwrap() error {
	return e.Err
}

func (e *HandlerError) LogAttrs() []slog.Attr {
	return []slog.Attr{
		slog.String("request_method", e.Method),
		slog.String("request_url", e.URL),
	}
}

type UnmarshalError struct {
	err error
}

func NewUnmarshalError(err error) *UnmarshalError {
	return &UnmarshalError{err: err}
}

func (e *UnmarshalError) Error() string   { return e.err.Error() }
func (e *UnmarshalError) StatusCode() int { return http.StatusBadRequest }

func NewParseError(err error) *UnmarshalError {
	return &UnmarshalError{err: err}
}

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

func (e *ValidationError) Error() string   { return e.err.Error() }
func (e *ValidationError) StatusCode() int { return http.StatusBadRequest }
