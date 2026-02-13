package eventvehicle

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for event vehicle data access
type Repository interface {
	Add(ctx context.Context, ev *EventVehicle) error
	Remove(ctx context.Context, eventID, vehicleID uuid.UUID) error
	UpdateQuantity(ctx context.Context, eventID, vehicleID uuid.UUID, quantity int) error
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventVehicleWithDetails, error)
	GetEventBrandID(ctx context.Context, eventID uuid.UUID) (uuid.UUID, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Add adds a vehicle to an event
func (r *PostgresRepository) Add(ctx context.Context, ev *EventVehicle) error {
	query := `
		INSERT INTO event_vehicles (id, event_id, vehicle_id, quantity, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query, ev.ID, ev.EventID, ev.VehicleID, ev.Quantity, ev.CreatedAt)
	return err
}

// Remove removes a vehicle from an event
func (r *PostgresRepository) Remove(ctx context.Context, eventID, vehicleID uuid.UUID) error {
	query := `DELETE FROM event_vehicles WHERE event_id = $1 AND vehicle_id = $2`

	result, err := r.pool.Exec(ctx, query, eventID, vehicleID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrVehicleNotFound
	}

	return nil
}

// UpdateQuantity updates the quantity of a vehicle in an event
func (r *PostgresRepository) UpdateQuantity(ctx context.Context, eventID, vehicleID uuid.UUID, quantity int) error {
	query := `
		UPDATE event_vehicles
		SET quantity = $3
		WHERE event_id = $1 AND vehicle_id = $2
	`

	result, err := r.pool.Exec(ctx, query, eventID, vehicleID, quantity)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrVehicleNotFound
	}

	return nil
}

// ListByEvent retrieves all vehicles for an event with details
func (r *PostgresRepository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventVehicleWithDetails, error) {
	query := `
		SELECT ev.id, ev.event_id, ev.vehicle_id, ev.quantity, ev.created_at,
		       v.model_name, v.brand_id, b.name as brand_name
		FROM event_vehicles ev
		JOIN vehicles v ON ev.vehicle_id = v.id
		JOIN brands b ON v.brand_id = b.id
		WHERE ev.event_id = $1
		ORDER BY b.name, v.model_name
	`

	rows, err := r.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []*EventVehicleWithDetails
	for rows.Next() {
		ev := &EventVehicleWithDetails{}
		err := rows.Scan(
			&ev.ID,
			&ev.EventID,
			&ev.VehicleID,
			&ev.Quantity,
			&ev.CreatedAt,
			&ev.ModelName,
			&ev.BrandID,
			&ev.BrandName,
		)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

// GetEventBrandID retrieves the brand ID of an event (used for validation)
func (r *PostgresRepository) GetEventBrandID(ctx context.Context, eventID uuid.UUID) (uuid.UUID, error) {
	query := `SELECT brand_id FROM events WHERE id = $1`

	var brandID uuid.UUID
	err := r.pool.QueryRow(ctx, query, eventID).Scan(&brandID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrEventNotFound
		}
		return uuid.Nil, err
	}

	return brandID, nil
}
