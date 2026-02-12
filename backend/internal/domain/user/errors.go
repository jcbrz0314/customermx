package user

import "errors"

var (
	// Validation errors
	ErrNameRequired        = errors.New("name is required")
	ErrEmailRequired       = errors.New("email is required")
	ErrPasswordRequired    = errors.New("password is required")
	ErrRoleRequired        = errors.New("role is required")
	ErrBrandIDRequired     = errors.New("brand_id is required for BRAND role users")
	ErrBrandIDNotAllowed   = errors.New("brand_id is not allowed for non-BRAND users")

	// Business logic errors
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserInactive        = errors.New("user account is inactive")
)
