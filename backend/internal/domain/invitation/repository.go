package invitation

import (
	"context"
	"errors"
	"fmt"

	"github.com/customermx/backend/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for invitation data access
type Repository interface {
	Create(ctx context.Context, invitation *Invitation) error
	GetByID(ctx context.Context, id uuid.UUID) (*Invitation, error)
	GetByToken(ctx context.Context, token string) (*Invitation, error)
	GetByEmail(ctx context.Context, email string) (*Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*Invitation, error)
	ListPending(ctx context.Context) ([]*Invitation, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates a new invitation
func (r *PostgresRepository) Create(ctx context.Context, invitation *Invitation) error {
	query := `
		INSERT INTO invitations (id, email, role, brand_id, token, expires_at, accepted, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(
		ctx, query,
		invitation.ID, invitation.Email, invitation.Role, invitation.BrandID,
		invitation.Token, invitation.ExpiresAt, invitation.Accepted, invitation.CreatedAt,
	).Scan(&invitation.ID, &invitation.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	return nil
}

// GetByID retrieves an invitation by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	query := `
		SELECT id, email, role, brand_id, token, expires_at, accepted, created_at
		FROM invitations
		WHERE id = $1
	`

	invitation := &Invitation{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&invitation.ID, &invitation.Email, &invitation.Role, &invitation.BrandID,
		&invitation.Token, &invitation.ExpiresAt, &invitation.Accepted, &invitation.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	return invitation, nil
}

// GetByToken retrieves an invitation by token
func (r *PostgresRepository) GetByToken(ctx context.Context, token string) (*Invitation, error) {
	query := `
		SELECT id, email, role, brand_id, token, expires_at, accepted, created_at
		FROM invitations
		WHERE token = $1
	`

	invitation := &Invitation{}
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&invitation.ID, &invitation.Email, &invitation.Role, &invitation.BrandID,
		&invitation.Token, &invitation.ExpiresAt, &invitation.Accepted, &invitation.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, fmt.Errorf("failed to get invitation by token: %w", err)
	}

	return invitation, nil
}

// GetByEmail retrieves an invitation by email
func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (*Invitation, error) {
	query := `
		SELECT id, email, role, brand_id, token, expires_at, accepted, created_at
		FROM invitations
		WHERE email = $1 AND accepted = false
		ORDER BY created_at DESC
		LIMIT 1
	`

	invitation := &Invitation{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&invitation.ID, &invitation.Email, &invitation.Role, &invitation.BrandID,
		&invitation.Token, &invitation.ExpiresAt, &invitation.Accepted, &invitation.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, fmt.Errorf("failed to get invitation by email: %w", err)
	}

	return invitation, nil
}

// Update updates an existing invitation
func (r *PostgresRepository) Update(ctx context.Context, invitation *Invitation) error {
	query := `
		UPDATE invitations
		SET accepted = $2, expires_at = $3
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, invitation.ID, invitation.Accepted, invitation.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrInvitationNotFound
	}

	return nil
}

// Delete deletes an invitation by ID
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM invitations WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrInvitationNotFound
	}

	return nil
}

// List retrieves all invitations
func (r *PostgresRepository) List(ctx context.Context) ([]*Invitation, error) {
	query := `
		SELECT id, email, role, brand_id, token, expires_at, accepted, created_at
		FROM invitations
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*Invitation
	for rows.Next() {
		invitation := &Invitation{}
		err := rows.Scan(
			&invitation.ID, &invitation.Email, &invitation.Role, &invitation.BrandID,
			&invitation.Token, &invitation.ExpiresAt, &invitation.Accepted, &invitation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

// ListPending retrieves all pending (not accepted) invitations
func (r *PostgresRepository) ListPending(ctx context.Context) ([]*Invitation, error) {
	query := `
		SELECT id, email, role, brand_id, token, expires_at, accepted, created_at
		FROM invitations
		WHERE accepted = false AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*Invitation
	for rows.Next() {
		invitation := &Invitation{}
		var role string
		err := rows.Scan(
			&invitation.ID, &invitation.Email, &role, &invitation.BrandID,
			&invitation.Token, &invitation.ExpiresAt, &invitation.Accepted, &invitation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitation.Role = user.Role(role)
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}
