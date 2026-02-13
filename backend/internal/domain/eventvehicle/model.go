package eventvehicle

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventVehicle represents a vehicle presented at an event
type EventVehicle struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	VehicleID uuid.UUID `json:"vehicle_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// EventVehicleWithDetails includes vehicle and brand information via JOIN
type EventVehicleWithDetails struct {
	EventVehicle
	ModelName string    `json:"model_name"`
	BrandID   uuid.UUID `json:"brand_id"`
	BrandName string    `json:"brand_name"`
}

// AddVehicleRequest represents the request to add a vehicle to an event
type AddVehicleRequest struct {
	EventID   string `json:"event_id"`
	VehicleID string `json:"vehicle_id"`
	Quantity  int    `json:"quantity"`
}

// Validate validates the add vehicle request
func (r *AddVehicleRequest) Validate() error {
	if strings.TrimSpace(r.EventID) == "" {
		return ErrEventIDRequired
	}
	if _, err := uuid.Parse(r.EventID); err != nil {
		return ErrEventIDRequired
	}
	if strings.TrimSpace(r.VehicleID) == "" {
		return ErrVehicleIDRequired
	}
	if _, err := uuid.Parse(r.VehicleID); err != nil {
		return ErrVehicleIDRequired
	}
	if r.Quantity < 1 {
		return ErrQuantityInvalid
	}
	return nil
}

// UpdateVehicleQuantityRequest represents the request to update vehicle quantity
type UpdateVehicleQuantityRequest struct {
	Quantity int `json:"quantity"`
}

// Validate validates the update quantity request
func (r *UpdateVehicleQuantityRequest) Validate() error {
	if r.Quantity < 1 {
		return ErrQuantityInvalid
	}
	return nil
}
