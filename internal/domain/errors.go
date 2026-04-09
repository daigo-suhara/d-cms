package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrValidation    = errors.New("validation error")
	ErrInvalidField  = errors.New("invalid field value")
	ErrHasEntries    = errors.New("content model has existing entries")
)
