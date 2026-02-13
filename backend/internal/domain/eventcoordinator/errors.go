package eventcoordinator

import "errors"

// Validation errors
var (
	ErrEventIDRequired = errors.New("event ID is required")
	ErrUserIDRequired  = errors.New("user ID is required")
)

// Business logic errors
var (
	ErrCoordinatorAlreadyExists = errors.New("coordinator already assigned to event")
	ErrUserNotCoordinator       = errors.New("user is not a coordinator")
	ErrCoordinatorNotFound      = errors.New("coordinator assignment not found")
	ErrEventNotFound            = errors.New("event not found")
	ErrUserNotFound             = errors.New("user not found")
)
