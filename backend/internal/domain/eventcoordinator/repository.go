package eventcoordinator

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for event coordinator data access
type Repository interface {
	Assign(ctx context.Context, ec *EventCoordinator) error
	Remove(ctx context.Context, eventID, userID uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventCoordinatorWithUser, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*EventCoordinator, error)
	IsCoordinatorAssigned(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Assign assigns a coordinator to an event
func (r *PostgresRepository) Assign(ctx context.Context, ec *EventCoordinator) error {
	query := `
		INSERT INTO event_coordinators (id, event_id, user_id, assigned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (event_id, user_id) DO NOTHING
	`

	_, err := r.pool.Exec(ctx, query, ec.ID, ec.EventID, ec.UserID, ec.AssignedAt)
	return err
}

// Remove removes a coordinator from an event
func (r *PostgresRepository) Remove(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `DELETE FROM event_coordinators WHERE event_id = $1 AND user_id = $2`

	result, err := r.pool.Exec(ctx, query, eventID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCoordinatorNotFound
	}

	return nil
}

// ListByEvent retrieves all coordinators for an event with user details
func (r *PostgresRepository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventCoordinatorWithUser, error) {
	query := `
		SELECT ec.id, ec.event_id, ec.user_id, ec.assigned_at, u.name, u.email
		FROM event_coordinators ec
		JOIN users u ON ec.user_id = u.id
		WHERE ec.event_id = $1
		ORDER BY ec.assigned_at
	`

	rows, err := r.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coordinators []*EventCoordinatorWithUser
	for rows.Next() {
		ec := &EventCoordinatorWithUser{}
		err := rows.Scan(
			&ec.ID,
			&ec.EventID,
			&ec.UserID,
			&ec.AssignedAt,
			&ec.UserName,
			&ec.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return coordinators, nil
}

// ListByUser retrieves all events assigned to a coordinator
func (r *PostgresRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*EventCoordinator, error) {
	query := `
		SELECT id, event_id, user_id, assigned_at
		FROM event_coordinators
		WHERE user_id = $1
		ORDER BY assigned_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coordinators []*EventCoordinator
	for rows.Next() {
		ec := &EventCoordinator{}
		err := rows.Scan(&ec.ID, &ec.EventID, &ec.UserID, &ec.AssignedAt)
		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return coordinators, nil
}

// IsCoordinatorAssigned checks if a user is assigned as coordinator to an event
func (r *PostgresRepository) IsCoordinatorAssigned(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM event_coordinators
			WHERE event_id = $1 AND user_id = $2
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, eventID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
