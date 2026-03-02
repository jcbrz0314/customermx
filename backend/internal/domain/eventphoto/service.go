package eventphoto

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StorageService abstracts S3 operations
type StorageService interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string, size int64) error
	Download(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}

// Service defines the interface for event photo business logic
type Service interface {
	Upload(ctx context.Context, eventID uuid.UUID, filename string, contentType string, size int64, body io.Reader) (*EventPhoto, error)
	List(ctx context.Context, eventID uuid.UUID) ([]*EventPhoto, error)
	GetForDownload(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID) (io.ReadCloser, string, error)
	Delete(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID) error
	Replace(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID, filename string, contentType string, size int64, body io.Reader) (*EventPhoto, error)
}

// PhotoService implements Service
type PhotoService struct {
	repo    Repository
	storage StorageService
}

// NewService creates a new PhotoService
func NewService(repo Repository, storage StorageService) *PhotoService {
	return &PhotoService{repo: repo, storage: storage}
}

// allowedImageTypes contains the accepted MIME types
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

const maxFileSizeBytes = 10 * 1024 * 1024 // 10 MB

func validateUpload(contentType string, size int64) error {
	// Normalize: "image/jpeg; charset=..." → "image/jpeg"
	ct := strings.Split(contentType, ";")[0]
	ct = strings.TrimSpace(strings.ToLower(ct))

	if !allowedImageTypes[ct] {
		return ErrInvalidContentType
	}
	if size > maxFileSizeBytes {
		return ErrFileTooLarge
	}
	return nil
}

// s3KeyForPhoto returns the S3 key for a given event + photo combination
func s3KeyForPhoto(eventID, photoID uuid.UUID) string {
	return fmt.Sprintf("events/%s/%s", eventID, photoID)
}

// photoURL returns the API proxy URL for a photo (relative to the /api/v1 prefix)
func photoURL(eventID, photoID uuid.UUID) string {
	return fmt.Sprintf("/events/%s/photos/%s", eventID, photoID)
}

// Upload stores a new photo for an event
func (s *PhotoService) Upload(ctx context.Context, eventID uuid.UUID, filename string, contentType string, size int64, body io.Reader) (*EventPhoto, error) {
	if err := validateUpload(contentType, size); err != nil {
		return nil, err
	}

	count, err := s.repo.CountByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if count >= MaxPhotosPerEvent {
		return nil, ErrMaxPhotosReached
	}

	photoID := uuid.New()
	key := s3KeyForPhoto(eventID, photoID)

	if err := s.storage.Upload(ctx, key, body, contentType, size); err != nil {
		return nil, err
	}

	photo := &EventPhoto{
		ID:          photoID,
		EventID:     eventID,
		S3Key:       key,
		Filename:    filename,
		ContentType: contentType,
		SizeBytes:   size,
		SortOrder:   count, // append at the end
		UploadedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, photo); err != nil {
		// Best-effort cleanup from S3
		_ = s.storage.Delete(ctx, key)
		return nil, err
	}

	photo.URL = photoURL(eventID, photoID)
	return photo, nil
}

// List returns all photos for an event with their proxy URLs
func (s *PhotoService) List(ctx context.Context, eventID uuid.UUID) ([]*EventPhoto, error) {
	photos, err := s.repo.ListByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	for _, p := range photos {
		p.URL = photoURL(p.EventID, p.ID)
	}
	return photos, nil
}

// GetForDownload returns an S3 ReadCloser for streaming to the client
func (s *PhotoService) GetForDownload(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID) (io.ReadCloser, string, error) {
	photo, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return nil, "", err
	}
	// Security: make sure the photo belongs to the requested event
	if photo.EventID != eventID {
		return nil, "", ErrPhotoNotFound
	}

	body, contentType, err := s.storage.Download(ctx, photo.S3Key)
	if err != nil {
		return nil, "", err
	}
	// Use stored content type if S3 doesn't return one
	if contentType == "" {
		contentType = photo.ContentType
	}
	return body, contentType, nil
}

// Delete removes a photo from S3 and the database
func (s *PhotoService) Delete(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID) error {
	photo, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}
	if photo.EventID != eventID {
		return ErrPhotoNotFound
	}

	// Delete from S3 first, then from DB
	if err := s.storage.Delete(ctx, photo.S3Key); err != nil {
		return err
	}
	return s.repo.Delete(ctx, photoID)
}

// Replace swaps the file behind an existing photo record
func (s *PhotoService) Replace(ctx context.Context, eventID uuid.UUID, photoID uuid.UUID, filename string, contentType string, size int64, body io.Reader) (*EventPhoto, error) {
	if err := validateUpload(contentType, size); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return nil, err
	}
	if existing.EventID != eventID {
		return nil, ErrPhotoNotFound
	}

	oldKey := existing.S3Key
	// Reuse the same key so the photo ID stays the same
	newKey := s3KeyForPhoto(eventID, photoID)

	if err := s.storage.Upload(ctx, newKey, body, contentType, size); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateS3Key(ctx, photoID, newKey, filename, contentType, size); err != nil {
		return nil, err
	}

	// Delete old object only if key changed (shouldn't happen with our current scheme)
	if oldKey != newKey {
		_ = s.storage.Delete(ctx, oldKey)
	}

	existing.Filename = filename
	existing.ContentType = contentType
	existing.SizeBytes = size
	existing.S3Key = newKey
	existing.URL = photoURL(eventID, photoID)
	return existing, nil
}
