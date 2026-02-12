package vehicle

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for vehicle data access
type Repository interface {
	Create(ctx context.Context, vehicle *Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error)
	GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*VehicleWithBrand, error)
	Update(ctx context.Context, vehicle *Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*VehicleWithBrand, error)
	ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Vehicle, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates a new vehicle
func (r *PostgresRepository) Create(ctx context.Context, vehicle *Vehicle) error {
	query := `
		INSERT INTO vehicles (id, brand_id, model_name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, vehicle.ID, vehicle.BrandID, vehicle.ModelName, vehicle.CreatedAt).
		Scan(&vehicle.ID, &vehicle.CreatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"vehicles_brand_id_model_name_key\" (SQLSTATE 23505)" {
			return ErrVehicleAlreadyExists
		}
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	return nil
}

// GetByID retrieves a vehicle by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error) {
	query := `
		SELECT id, brand_id, model_name, created_at
		FROM vehicles
		WHERE id = $1
	`

	vehicle := &Vehicle{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&vehicle.ID, &vehicle.BrandID, &vehicle.ModelName, &vehicle.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVehicleNotFound
		}
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return vehicle, nil
}

// GetByIDWithBrand retrieves a vehicle by ID with brand information
func (r *PostgresRepository) GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*VehicleWithBrand, error) {
	query := `
		SELECT v.id, v.brand_id, b.name as brand_name, v.model_name, v.created_at
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.id
		WHERE v.id = $1
	`

	vehicle := &VehicleWithBrand{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&vehicle.ID, &vehicle.BrandID, &vehicle.BrandName, &vehicle.ModelName, &vehicle.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVehicleNotFound
		}
		return nil, fmt.Errorf("failed to get vehicle with brand: %w", err)
	}

	return vehicle, nil
}

// Update updates an existing vehicle
func (r *PostgresRepository) Update(ctx context.Context, vehicle *Vehicle) error {
	query := `
		UPDATE vehicles
		SET model_name = $2
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, vehicle.ID, vehicle.ModelName)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrVehicleNotFound
	}

	return nil
}

// Delete deletes a vehicle by ID
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vehicles WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrVehicleNotFound
	}

	return nil
}

// List retrieves all vehicles with brand information
func (r *PostgresRepository) List(ctx context.Context) ([]*VehicleWithBrand, error) {
	query := `
		SELECT v.id, v.brand_id, b.name as brand_name, v.model_name, v.created_at
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.id
		ORDER BY b.name, v.model_name ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*VehicleWithBrand
	for rows.Next() {
		vehicle := &VehicleWithBrand{}
		err := rows.Scan(
			&vehicle.ID, &vehicle.BrandID, &vehicle.BrandName,
			&vehicle.ModelName, &vehicle.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// ListByBrandID retrieves all vehicles for a specific brand
func (r *PostgresRepository) ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Vehicle, error) {
	query := `
		SELECT id, brand_id, model_name, created_at
		FROM vehicles
		WHERE brand_id = $1
		ORDER BY model_name ASC
	`

	rows, err := r.pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles by brand: %w", err)
	}
	defer rows.Close()

	var vehicles []*Vehicle
	for rows.Next() {
		vehicle := &Vehicle{}
		err := rows.Scan(
			&vehicle.ID, &vehicle.BrandID, &vehicle.ModelName, &vehicle.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}
