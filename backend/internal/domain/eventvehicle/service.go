package eventvehicle

import (
	"context"
	"time"

	"github.com/customermx/backend/internal/domain/vehicle"
	"github.com/google/uuid"
)

// Service defines the interface for event vehicle business logic
type Service interface {
	AddVehicle(ctx context.Context, req *AddVehicleRequest) (*EventVehicle, error)
	RemoveVehicle(ctx context.Context, eventID, vehicleID uuid.UUID) error
	UpdateQuantity(ctx context.Context, eventID, vehicleID uuid.UUID, req *UpdateVehicleQuantityRequest) error
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventVehicleWithDetails, error)
}

// EventVehicleService implements the Service interface
type EventVehicleService struct {
	repo        Repository
	vehicleRepo vehicle.Repository
}

// NewService creates a new EventVehicleService
func NewService(repo Repository, vehicleRepo vehicle.Repository) *EventVehicleService {
	return &EventVehicleService{
		repo:        repo,
		vehicleRepo: vehicleRepo,
	}
}

// AddVehicle adds a vehicle to an event with brand validation
func (s *EventVehicleService) AddVehicle(ctx context.Context, req *AddVehicleRequest) (*EventVehicle, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, ErrEventIDRequired
	}

	vehicleID, err := uuid.Parse(req.VehicleID)
	if err != nil {
		return nil, ErrVehicleIDRequired
	}

	// CRITICAL VALIDATION: Get event's brand ID
	eventBrandID, err := s.repo.GetEventBrandID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Get vehicle to check its brand
	veh, err := s.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		if err == vehicle.ErrVehicleNotFound {
			return nil, ErrVehicleNotFound
		}
		return nil, err
	}

	// CRITICAL: Validate that vehicle brand matches event brand
	if veh.BrandID != eventBrandID {
		return nil, ErrBrandMismatch
	}

	// Create event vehicle
	ev := &EventVehicle{
		ID:        uuid.New(),
		EventID:   eventID,
		VehicleID: vehicleID,
		Quantity:  req.Quantity,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Add(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

// RemoveVehicle removes a vehicle from an event
func (s *EventVehicleService) RemoveVehicle(ctx context.Context, eventID, vehicleID uuid.UUID) error {
	return s.repo.Remove(ctx, eventID, vehicleID)
}

// UpdateQuantity updates the quantity of a vehicle in an event
func (s *EventVehicleService) UpdateQuantity(ctx context.Context, eventID, vehicleID uuid.UUID, req *UpdateVehicleQuantityRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateQuantity(ctx, eventID, vehicleID, req.Quantity)
}

// ListByEvent retrieves all vehicles for an event
func (s *EventVehicleService) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventVehicleWithDetails, error) {
	return s.repo.ListByEvent(ctx, eventID)
}
