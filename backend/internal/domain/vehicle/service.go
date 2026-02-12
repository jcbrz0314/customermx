package vehicle

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service defines the vehicle business logic interface
type Service interface {
	CreateVehicle(ctx context.Context, req *CreateVehicleRequest) (*Vehicle, error)
	GetVehicle(ctx context.Context, id uuid.UUID) (*VehicleWithBrand, error)
	UpdateVehicle(ctx context.Context, id uuid.UUID, req *UpdateVehicleRequest) (*Vehicle, error)
	DeleteVehicle(ctx context.Context, id uuid.UUID) error
	ListVehicles(ctx context.Context) ([]*VehicleWithBrand, error)
	ListVehiclesByBrand(ctx context.Context, brandID uuid.UUID) ([]*Vehicle, error)
}

// VehicleService implements the Service interface
type VehicleService struct {
	repo Repository
}

// NewService creates a new VehicleService
func NewService(repo Repository) *VehicleService {
	return &VehicleService{repo: repo}
}

// CreateVehicle creates a new vehicle
func (s *VehicleService) CreateVehicle(ctx context.Context, req *CreateVehicleRequest) (*Vehicle, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	vehicle := &Vehicle{
		ID:        uuid.New(),
		BrandID:   req.BrandID,
		ModelName: req.ModelName,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// GetVehicle retrieves a vehicle by ID with brand information
func (s *VehicleService) GetVehicle(ctx context.Context, id uuid.UUID) (*VehicleWithBrand, error) {
	return s.repo.GetByIDWithBrand(ctx, id)
}

// UpdateVehicle updates an existing vehicle
func (s *VehicleService) UpdateVehicle(ctx context.Context, id uuid.UUID, req *UpdateVehicleRequest) (*Vehicle, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	vehicle, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicle.ModelName = req.ModelName

	if err := s.repo.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// DeleteVehicle deletes a vehicle
func (s *VehicleService) DeleteVehicle(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ListVehicles retrieves all vehicles with brand information
func (s *VehicleService) ListVehicles(ctx context.Context) ([]*VehicleWithBrand, error) {
	return s.repo.List(ctx)
}

// ListVehiclesByBrand retrieves all vehicles for a specific brand
func (s *VehicleService) ListVehiclesByBrand(ctx context.Context, brandID uuid.UUID) ([]*Vehicle, error) {
	return s.repo.ListByBrandID(ctx, brandID)
}
