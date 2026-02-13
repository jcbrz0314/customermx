package event

import "errors"

// Validation errors
var (
	ErrBrandIDRequired   = errors.New("brand ID is required")
	ErrNameRequired      = errors.New("event name is required")
	ErrEventTypeRequired = errors.New("event type is required")
	ErrStartDateRequired = errors.New("start date is required")
	ErrStartDateInvalid  = errors.New("start date must be in YYYY-MM-DD format")
	ErrYearInvalid       = errors.New("year must be between 2000 and 2100")
	ErrDurationInvalid   = errors.New("duration must be at least 1 day")
	ErrInvalidStatus     = errors.New("invalid event status")
)

// Business logic errors
var (
	ErrEventNotFound          = errors.New("event not found")
	ErrBrandNotFound          = errors.New("brand not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
