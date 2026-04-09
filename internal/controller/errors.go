package controller

import (
	"errors"
	"net/http"

	"github.com/daigo-suhara/d-cms/internal/domain"
)

func httpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrInvalidField):
		return http.StatusUnprocessableEntity
	case errors.Is(err, domain.ErrHasEntries):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
