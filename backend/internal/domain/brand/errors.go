package brand

import "errors"

var (
	// Validation errors
	ErrNameRequired = errors.New("brand name is required")

	// Business logic errors
	ErrBrandNotFound      = errors.New("brand not found")
	ErrBrandAlreadyExists = errors.New("brand with this name already exists")
)
