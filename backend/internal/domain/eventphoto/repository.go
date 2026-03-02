package eventphoto

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for event photo data access
type Repository interface {
	Create(ctx context.Context, photo *EventPhoto) error
	GetByID(ctx context.Context, id uuid.UUID) (*EventPhoto, error)
	ListByEventID(ctx context.Context, eventID uuid.UUID) ([]*EventPhoto, error)
	CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateS3Key(ctx context.Context, id uuid.UUID, s3Key string, filename string, contentType string, sizeBytes int64) error
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgresRepository
func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// Create inserts a new event photo record
func (r *PostgresRepository) Create(ctx context.Context, photo *EventPhoto) error {
	query := `
		INSERT INTO event_photos (id, event_id, s3_key, filename, content_type, size_bytes, sort_order, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		photo.ID,
		photo.EventID,
		photo.S3Key,
		photo.Filename,
		photo.ContentType,
		photo.SizeBytes,
		photo.SortOrder,
		photo.UploadedAt,
	)
	return err
}

// GetByID retrieves a photo by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*EventPhoto, error) {
	query := `
		SELECT id, event_id, s3_key, filename, content_type, size_bytes, sort_order, uploaded_at
		FROM event_photos
		WHERE id = $1
	`
	photo := &EventPhoto{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&photo.ID,
		&photo.EventID,
		&photo.S3Key,
		&photo.Filename,
		&photo.ContentType,
		&photo.SizeBytes,
		&photo.SortOrder,
		&photo.UploadedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPhotoNotFound
		}
		return nil, err
	}
	return photo, nil
}

// ListByEventID retrieves all photos for an event, ordered by sort_order
func (r *PostgresRepository) ListByEventID(ctx context.Context, eventID uuid.UUID) ([]*EventPhoto, error) {
	query := `
		SELECT id, event_id, s3_key, filename, content_type, size_bytes, sort_order, uploaded_at
		FROM event_photos
		WHERE event_id = $1
		ORDER BY sort_order ASC, uploaded_at ASC
	`
	rows, err := r.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*EventPhoto
	for rows.Next() {
		photo := &EventPhoto{}
		if err := rows.Scan(
			&photo.ID,
			&photo.EventID,
			&photo.S3Key,
			&photo.Filename,
			&photo.ContentType,
			&photo.SizeBytes,
			&photo.SortOrder,
			&photo.UploadedAt,
		); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return photos, nil
}

// CountByEventID returns the number of photos for an event
func (r *PostgresRepository) CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM event_photos WHERE event_id = $1`, eventID).Scan(&count)
	return count, err
}

// Delete removes a photo record
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM event_photos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrPhotoNotFound
	}
	return nil
}

// UpdateS3Key replaces the S3 key and metadata of a photo (used when replacing a photo)
func (r *PostgresRepository) UpdateS3Key(ctx context.Context, id uuid.UUID, s3Key string, filename string, contentType string, sizeBytes int64) error {
	query := `
		UPDATE event_photos
		SET s3_key = $2, filename = $3, content_type = $4, size_bytes = $5, uploaded_at = now()
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, id, s3Key, filename, contentType, sizeBytes)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrPhotoNotFound
	}
	return nil
}
