package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for user data access
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*User, error)
	ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*User, error)
	ListByRole(ctx context.Context, role Role) ([]*User, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates a new user
func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(
		ctx, query,
		user.ID, user.Name, user.Email, user.PasswordHash,
		user.Role, user.BrandID, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	var roleStr string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&roleStr, &user.BrandID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Role = Role(roleStr)
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	var roleStr string
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&roleStr, &user.BrandID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Role = Role(roleStr)
	return user, nil
}

// Update updates an existing user
func (r *PostgresRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, is_active = $4, updated_at = $5
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(
		ctx, query,
		user.ID, user.Name, user.Email, user.IsActive, user.UpdatedAt,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by ID
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// List retrieves all users
func (r *PostgresRepository) List(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		var roleStr string
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&roleStr, &user.BrandID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Role = Role(roleStr)
		users = append(users, user)
	}

	return users, nil
}

// ListByBrandID retrieves all users for a specific brand
func (r *PostgresRepository) ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at
		FROM users
		WHERE brand_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by brand: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		var roleStr string
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&roleStr, &user.BrandID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Role = Role(roleStr)
		users = append(users, user)
	}

	return users, nil
}

// ListByRole retrieves all users with a specific role
func (r *PostgresRepository) ListByRole(ctx context.Context, role Role) ([]*User, error) {
	query := `
		SELECT id, name, email, password_hash, role, brand_id, is_active, created_at, updated_at
		FROM users
		WHERE role = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by role: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		var roleStr string
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordHash,
			&roleStr, &user.BrandID, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.Role = Role(roleStr)
		users = append(users, user)
	}

	return users, nil
}
