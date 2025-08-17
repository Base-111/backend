package api

import (
	"strings"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type FieldErrors = map[string][]string

type ValidationError struct {
	fieldErrors FieldErrors
}

func NewValidationError() *ValidationError {
	return &ValidationError{fieldErrors: make(FieldErrors)}
}

func (e *ValidationError) Add(field string, err string) {
	e.fieldErrors[field] = append(e.fieldErrors[field], err)
}

func (e *ValidationError) Error() string {
	sb := new(strings.Builder)
	sb.WriteString("validation error:\n")

	for k, v := range e.fieldErrors {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(strings.Join(v, ","))
		sb.WriteString("\n")
	}

	return sb.String()
}

func (e *ValidationError) GetValidationErrors() FieldErrors {
	return e.fieldErrors
}
