package eventreport

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for event report data access
type Repository interface {
	Create(ctx context.Context, report *EventReport) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) (*EventReport, error)
	Update(ctx context.Context, report *EventReport) error
	MarkAsCompleted(ctx context.Context, eventID uuid.UUID, completed bool) error
	Delete(ctx context.Context, eventID uuid.UUID) error
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create creates or updates an event report (UPSERT)
func (r *PostgresRepository) Create(ctx context.Context, report *EventReport) error {
	query := `
		INSERT INTO event_reports (
			id, event_id, hostess_count, setup_vendor, has_promotional,
			attendees, activities_count, leads_collected, prospects,
			dealer_rating, comments, completed, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (event_id) DO UPDATE SET
			hostess_count = COALESCE(EXCLUDED.hostess_count, event_reports.hostess_count),
			setup_vendor = COALESCE(EXCLUDED.setup_vendor, event_reports.setup_vendor),
			has_promotional = COALESCE(EXCLUDED.has_promotional, event_reports.has_promotional),
			attendees = COALESCE(EXCLUDED.attendees, event_reports.attendees),
			activities_count = COALESCE(EXCLUDED.activities_count, event_reports.activities_count),
			leads_collected = COALESCE(EXCLUDED.leads_collected, event_reports.leads_collected),
			prospects = COALESCE(EXCLUDED.prospects, event_reports.prospects),
			dealer_rating = COALESCE(EXCLUDED.dealer_rating, event_reports.dealer_rating),
			comments = COALESCE(EXCLUDED.comments, event_reports.comments),
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		report.ID,
		report.EventID,
		report.HostessCount,
		report.SetupVendor,
		report.HasPromotional,
		report.Attendees,
		report.ActivitiesCount,
		report.LeadsCollected,
		report.Prospects,
		report.DealerRating,
		report.Comments,
		report.Completed,
		report.CreatedAt,
		report.UpdatedAt,
	).Scan(&report.ID)

	return err
}

// GetByEventID retrieves an event report by event ID
func (r *PostgresRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) (*EventReport, error) {
	query := `
		SELECT id, event_id, hostess_count, setup_vendor, has_promotional,
		       attendees, activities_count, leads_collected, prospects,
		       dealer_rating, comments, completed, created_at, updated_at
		FROM event_reports
		WHERE event_id = $1
	`

	report := &EventReport{}
	err := r.pool.QueryRow(ctx, query, eventID).Scan(
		&report.ID,
		&report.EventID,
		&report.HostessCount,
		&report.SetupVendor,
		&report.HasPromotional,
		&report.Attendees,
		&report.ActivitiesCount,
		&report.LeadsCollected,
		&report.Prospects,
		&report.DealerRating,
		&report.Comments,
		&report.Completed,
		&report.CreatedAt,
		&report.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrReportNotFound
		}
		return nil, err
	}

	return report, nil
}

// Update updates an existing event report
func (r *PostgresRepository) Update(ctx context.Context, report *EventReport) error {
	query := `
		UPDATE event_reports
		SET hostess_count = $2, setup_vendor = $3, has_promotional = $4,
		    attendees = $5, activities_count = $6, leads_collected = $7,
		    prospects = $8, dealer_rating = $9, comments = $10,
		    updated_at = $11
		WHERE event_id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		report.EventID,
		report.HostessCount,
		report.SetupVendor,
		report.HasPromotional,
		report.Attendees,
		report.ActivitiesCount,
		report.LeadsCollected,
		report.Prospects,
		report.DealerRating,
		report.Comments,
		time.Now(),
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrReportNotFound
	}

	return nil
}

// MarkAsCompleted marks a report as completed or not completed
func (r *PostgresRepository) MarkAsCompleted(ctx context.Context, eventID uuid.UUID, completed bool) error {
	query := `
		UPDATE event_reports
		SET completed = $2, updated_at = $3
		WHERE event_id = $1
	`

	result, err := r.pool.Exec(ctx, query, eventID, completed, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrReportNotFound
	}

	return nil
}

// Delete deletes an event report
func (r *PostgresRepository) Delete(ctx context.Context, eventID uuid.UUID) error {
	query := `DELETE FROM event_reports WHERE event_id = $1`

	result, err := r.pool.Exec(ctx, query, eventID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrReportNotFound
	}

	return nil
}
