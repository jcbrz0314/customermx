package vehicle

import "errors"

var (
	// Validation errors
	ErrBrandIDRequired    = errors.New("brand_id is required")
	ErrModelNameRequired  = errors.New("model_name is required")

	// Business logic errors
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists for this brand")
)
