package event

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for event data access
type Repository interface {
	Create(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*EventWithBrand, error)
	Update(ctx context.Context, event *Event) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status EventStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters EventFilters) ([]*EventWithBrand, error)
	ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Event, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates a new event
func (r *PostgresRepository) Create(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO events (
			id, brand_id, event_type, organizer, name, start_date, year,
			duration_days, state, city, venue, dealer, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.pool.Exec(ctx, query,
		event.ID,
		event.BrandID,
		event.EventType,
		event.Organizer,
		event.Name,
		event.StartDate,
		event.Year,
		event.DurationDays,
		event.State,
		event.City,
		event.Venue,
		event.Dealer,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	)

	return err
}

// GetByID retrieves an event by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Event, error) {
	query := `
		SELECT id, brand_id, event_type, organizer, name, start_date, year,
		       duration_days, state, city, venue, dealer, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	event := &Event{}
	var statusStr string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.BrandID,
		&event.EventType,
		&event.Organizer,
		&event.Name,
		&event.StartDate,
		&event.Year,
		&event.DurationDays,
		&event.State,
		&event.City,
		&event.Venue,
		&event.Dealer,
		&statusStr,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	event.Status = EventStatus(statusStr)
	return event, nil
}

// GetByIDWithBrand retrieves an event by ID with brand name
func (r *PostgresRepository) GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*EventWithBrand, error) {
	query := `
		SELECT e.id, e.brand_id, e.event_type, e.organizer, e.name, e.start_date, e.year,
		       e.duration_days, e.state, e.city, e.venue, e.dealer, e.status,
		       e.created_at, e.updated_at, b.name as brand_name
		FROM events e
		JOIN brands b ON e.brand_id = b.id
		WHERE e.id = $1
	`

	eventWithBrand := &EventWithBrand{}
	var statusStr string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&eventWithBrand.ID,
		&eventWithBrand.BrandID,
		&eventWithBrand.EventType,
		&eventWithBrand.Organizer,
		&eventWithBrand.Name,
		&eventWithBrand.StartDate,
		&eventWithBrand.Year,
		&eventWithBrand.DurationDays,
		&eventWithBrand.State,
		&eventWithBrand.City,
		&eventWithBrand.Venue,
		&eventWithBrand.Dealer,
		&statusStr,
		&eventWithBrand.CreatedAt,
		&eventWithBrand.UpdatedAt,
		&eventWithBrand.BrandName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	eventWithBrand.Status = EventStatus(statusStr)
	return eventWithBrand, nil
}

// Update updates an existing event
func (r *PostgresRepository) Update(ctx context.Context, event *Event) error {
	query := `
		UPDATE events
		SET event_type = $2, organizer = $3, name = $4, start_date = $5, year = $6,
		    duration_days = $7, state = $8, city = $9, venue = $10, dealer = $11,
		    updated_at = $12
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		event.ID,
		event.EventType,
		event.Organizer,
		event.Name,
		event.StartDate,
		event.Year,
		event.DurationDays,
		event.State,
		event.City,
		event.Venue,
		event.Dealer,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEventNotFound
	}

	return nil
}

// UpdateStatus updates the status of an event
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status EventStatus) error {
	query := `
		UPDATE events
		SET status = $2, updated_at = $3
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id, status, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEventNotFound
	}

	return nil
}

// Delete deletes an event
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEventNotFound
	}

	return nil
}

// List retrieves events with optional filters
func (r *PostgresRepository) List(ctx context.Context, filters EventFilters) ([]*EventWithBrand, error) {
	query := `
		SELECT e.id, e.brand_id, e.event_type, e.organizer, e.name, e.start_date, e.year,
		       e.duration_days, e.state, e.city, e.venue, e.dealer, e.status,
		       e.created_at, e.updated_at, b.name as brand_name
		FROM events e
		JOIN brands b ON e.brand_id = b.id
	`

	// If filtering by coordinator, JOIN with event_coordinators
	if filters.CoordinatorID != nil {
		query += ` JOIN event_coordinators ec ON ec.event_id = e.id AND ec.user_id = $5`
	}

	query += `
		WHERE 1=1
		  AND ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND ($3::event_status IS NULL OR e.status = $3)
		  AND ($4::text IS NULL OR e.state = $4)
		ORDER BY e.start_date DESC
	`

	var rows pgx.Rows
	var err error
	if filters.CoordinatorID != nil {
		rows, err = r.pool.Query(ctx, query, filters.BrandID, filters.Year, filters.Status, filters.State, filters.CoordinatorID)
	} else {
		rows, err = r.pool.Query(ctx, query, filters.BrandID, filters.Year, filters.Status, filters.State)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*EventWithBrand
	for rows.Next() {
		eventWithBrand := &EventWithBrand{}
		var statusStr string

		err := rows.Scan(
			&eventWithBrand.ID,
			&eventWithBrand.BrandID,
			&eventWithBrand.EventType,
			&eventWithBrand.Organizer,
			&eventWithBrand.Name,
			&eventWithBrand.StartDate,
			&eventWithBrand.Year,
			&eventWithBrand.DurationDays,
			&eventWithBrand.State,
			&eventWithBrand.City,
			&eventWithBrand.Venue,
			&eventWithBrand.Dealer,
			&statusStr,
			&eventWithBrand.CreatedAt,
			&eventWithBrand.UpdatedAt,
			&eventWithBrand.BrandName,
		)
		if err != nil {
			return nil, err
		}

		eventWithBrand.Status = EventStatus(statusStr)
		events = append(events, eventWithBrand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// ListByBrandID retrieves all events for a specific brand
func (r *PostgresRepository) ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Event, error) {
	query := `
		SELECT id, brand_id, event_type, organizer, name, start_date, year,
		       duration_days, state, city, venue, dealer, status, created_at, updated_at
		FROM events
		WHERE brand_id = $1
		ORDER BY start_date DESC
	`

	rows, err := r.pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event := &Event{}
		var statusStr string

		err := rows.Scan(
			&event.ID,
			&event.BrandID,
			&event.EventType,
			&event.Organizer,
			&event.Name,
			&event.StartDate,
			&event.Year,
			&event.DurationDays,
			&event.State,
			&event.City,
			&event.Venue,
			&event.Dealer,
			&statusStr,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		event.Status = EventStatus(statusStr)
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
