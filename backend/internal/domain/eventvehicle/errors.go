package eventvehicle

import "errors"

// Validation errors
var (
	ErrEventIDRequired   = errors.New("event ID is required")
	ErrVehicleIDRequired = errors.New("vehicle ID is required")
	ErrQuantityInvalid   = errors.New("quantity must be at least 1")
)

// Business logic errors
var (
	ErrVehicleAlreadyAdded = errors.New("vehicle already added to event")
	ErrEventNotFound       = errors.New("event not found")
	ErrVehicleNotFound     = errors.New("vehicle not found")
	ErrBrandMismatch       = errors.New("vehicle brand does not match event brand")
)
