package eventreport

import (
	"context"
	"time"

	"github.com/customermx/backend/internal/domain/event"
	"github.com/google/uuid"
)

// Service defines the interface for event report business logic
type Service interface {
	CreateOrUpdate(ctx context.Context, eventID uuid.UUID, req *CreateOrUpdateReportRequest) (*EventReport, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) (*EventReport, error)
	MarkAsCompleted(ctx context.Context, eventID uuid.UUID, req *CompleteReportRequest) (*EventReport, error)
	Delete(ctx context.Context, eventID uuid.UUID) error
}

// EventReportService implements the Service interface
type EventReportService struct {
	repo      Repository
	eventRepo event.Repository
}

// NewService creates a new EventReportService
func NewService(repo Repository, eventRepo event.Repository) *EventReportService {
	return &EventReportService{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

// CreateOrUpdate creates or updates an event report (UPSERT)
func (s *EventReportService) CreateOrUpdate(ctx context.Context, eventID uuid.UUID, req *CreateOrUpdateReportRequest) (*EventReport, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify event exists
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if err == event.ErrEventNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	// Create or update report
	report := &EventReport{
		ID:              uuid.New(),
		EventID:         eventID,
		HostessCount:    req.HostessCount,
		SetupVendor:     req.SetupVendor,
		HasPromotional:  req.HasPromotional,
		Attendees:       req.Attendees,
		ActivitiesCount: req.ActivitiesCount,
		LeadsCollected:  req.LeadsCollected,
		Prospects:       req.Prospects,
		DealerRating:    req.DealerRating,
		Comments:        req.Comments,
		Completed:       false, // Default to not completed
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}

	// Fetch the final report to return (in case of UPSERT)
	return s.repo.GetByEventID(ctx, eventID)
}

// GetByEventID retrieves an event report by event ID
func (s *EventReportService) GetByEventID(ctx context.Context, eventID uuid.UUID) (*EventReport, error) {
	return s.repo.GetByEventID(ctx, eventID)
}

// MarkAsCompleted marks a report as completed
func (s *EventReportService) MarkAsCompleted(ctx context.Context, eventID uuid.UUID, req *CompleteReportRequest) (*EventReport, error) {
	// If marking as completed, verify event is in COMPLETED or CLOSED status
	if req.Completed {
		evt, err := s.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			if err == event.ErrEventNotFound {
				return nil, ErrEventNotFound
			}
			return nil, err
		}

		if evt.Status != event.StatusCompleted && evt.Status != event.StatusClosed {
			return nil, ErrEventNotCompleted
		}
	}

	// Update completed status
	if err := s.repo.MarkAsCompleted(ctx, eventID, req.Completed); err != nil {
		return nil, err
	}

	return s.repo.GetByEventID(ctx, eventID)
}

// Delete deletes an event report
func (s *EventReportService) Delete(ctx context.Context, eventID uuid.UUID) error {
	return s.repo.Delete(ctx, eventID)
}
