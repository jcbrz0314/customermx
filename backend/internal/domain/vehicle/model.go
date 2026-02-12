package vehicle

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle represents a vehicle model
type Vehicle struct {
	ID        uuid.UUID `json:"id"`
	BrandID   uuid.UUID `json:"brand_id"`
	ModelName string    `json:"model_name"`
	CreatedAt time.Time `json:"created_at"`
}

// VehicleWithBrand includes brand information
type VehicleWithBrand struct {
	ID        uuid.UUID `json:"id"`
	BrandID   uuid.UUID `json:"brand_id"`
	BrandName string    `json:"brand_name"`
	ModelName string    `json:"model_name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateVehicleRequest represents the request to create a new vehicle
type CreateVehicleRequest struct {
	BrandID   uuid.UUID `json:"brand_id"`
	ModelName string    `json:"model_name"`
}

// UpdateVehicleRequest represents the request to update a vehicle
type UpdateVehicleRequest struct {
	ModelName string `json:"model_name"`
}

// Validate validates the CreateVehicleRequest
func (r *CreateVehicleRequest) Validate() error {
	if r.BrandID == uuid.Nil {
		return ErrBrandIDRequired
	}
	if r.ModelName == "" {
		return ErrModelNameRequired
	}
	return nil
}

// Validate validates the UpdateVehicleRequest
func (r *UpdateVehicleRequest) Validate() error {
	if r.ModelName == "" {
		return ErrModelNameRequired
	}
	return nil
}
