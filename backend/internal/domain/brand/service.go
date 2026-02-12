package brand

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service defines the brand business logic interface
type Service interface {
	CreateBrand(ctx context.Context, req *CreateBrandRequest) (*Brand, error)
	GetBrand(ctx context.Context, id uuid.UUID) (*Brand, error)
	GetBrandByName(ctx context.Context, name string) (*Brand, error)
	UpdateBrand(ctx context.Context, id uuid.UUID, req *UpdateBrandRequest) (*Brand, error)
	DeleteBrand(ctx context.Context, id uuid.UUID) error
	ListBrands(ctx context.Context) ([]*Brand, error)
}

// BrandService implements the Service interface
type BrandService struct {
	repo Repository
}

// NewService creates a new BrandService
func NewService(repo Repository) *BrandService {
	return &BrandService{repo: repo}
}

// CreateBrand creates a new brand
func (s *BrandService) CreateBrand(ctx context.Context, req *CreateBrandRequest) (*Brand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	brand := &Brand{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// GetBrand retrieves a brand by ID
func (s *BrandService) GetBrand(ctx context.Context, id uuid.UUID) (*Brand, error) {
	return s.repo.GetByID(ctx, id)
}

// GetBrandByName retrieves a brand by name
func (s *BrandService) GetBrandByName(ctx context.Context, name string) (*Brand, error) {
	return s.repo.GetByName(ctx, name)
}

// UpdateBrand updates an existing brand
func (s *BrandService) UpdateBrand(ctx context.Context, id uuid.UUID, req *UpdateBrandRequest) (*Brand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	brand.Name = req.Name

	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// DeleteBrand deletes a brand
func (s *BrandService) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ListBrands retrieves all brands
func (s *BrandService) ListBrands(ctx context.Context) ([]*Brand, error) {
	return s.repo.List(ctx)
}
