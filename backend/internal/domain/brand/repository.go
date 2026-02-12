package brand

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for brand data access
type Repository interface {
	Create(ctx context.Context, brand *Brand) error
	GetByID(ctx context.Context, id uuid.UUID) (*Brand, error)
	GetByName(ctx context.Context, name string) (*Brand, error)
	Update(ctx context.Context, brand *Brand) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*Brand, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates a new brand
func (r *PostgresRepository) Create(ctx context.Context, brand *Brand) error {
	query := `
		INSERT INTO brands (id, name, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, brand.ID, brand.Name, brand.CreatedAt).
		Scan(&brand.ID, &brand.CreatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"brands_name_key\" (SQLSTATE 23505)" {
			return ErrBrandAlreadyExists
		}
		return fmt.Errorf("failed to create brand: %w", err)
	}

	return nil
}

// GetByID retrieves a brand by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Brand, error) {
	query := `SELECT id, name, created_at FROM brands WHERE id = $1`

	brand := &Brand{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&brand.ID, &brand.Name, &brand.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBrandNotFound
		}
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	return brand, nil
}

// GetByName retrieves a brand by name
func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*Brand, error) {
	query := `SELECT id, name, created_at FROM brands WHERE name = $1`

	brand := &Brand{}
	err := r.pool.QueryRow(ctx, query, name).Scan(&brand.ID, &brand.Name, &brand.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBrandNotFound
		}
		return nil, fmt.Errorf("failed to get brand by name: %w", err)
	}

	return brand, nil
}

// Update updates an existing brand
func (r *PostgresRepository) Update(ctx context.Context, brand *Brand) error {
	query := `UPDATE brands SET name = $2 WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, brand.ID, brand.Name)
	if err != nil {
		return fmt.Errorf("failed to update brand: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrBrandNotFound
	}

	return nil
}

// Delete deletes a brand by ID
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM brands WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete brand: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrBrandNotFound
	}

	return nil
}

// List retrieves all brands
func (r *PostgresRepository) List(ctx context.Context) ([]*Brand, error) {
	query := `SELECT id, name, created_at FROM brands ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list brands: %w", err)
	}
	defer rows.Close()

	var brands []*Brand
	for rows.Next() {
		brand := &Brand{}
		err := rows.Scan(&brand.ID, &brand.Name, &brand.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}
