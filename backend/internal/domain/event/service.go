package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for event business logic
type Service interface {
	Create(ctx context.Context, req *CreateEventRequest) (*Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*EventWithBrand, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateEventRequest) (*Event, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, req *ChangeStatusRequest) (*Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters EventFilters) ([]*EventWithBrand, error)
	ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Event, error)
}

// EventService implements the Service interface
type EventService struct {
	repo Repository
}

// NewService creates a new EventService
func NewService(repo Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

// Create creates a new event
func (s *EventService) Create(ctx context.Context, req *CreateEventRequest) (*Event, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	brandID, err := uuid.Parse(req.BrandID)
	if err != nil {
		return nil, ErrBrandIDRequired
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrStartDateInvalid
	}

	event := &Event{
		ID:           uuid.New(),
		BrandID:      brandID,
		EventType:    req.EventType,
		Organizer:    req.Organizer,
		Name:         req.Name,
		StartDate:    startDate,
		Year:         req.Year,
		DurationDays: req.DurationDays,
		State:        req.State,
		City:         req.City,
		Venue:        req.Venue,
		Dealer:       req.Dealer,
		Status:       StatusPlanned, // Default status
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// GetByID retrieves an event by ID
func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (*Event, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByIDWithBrand retrieves an event by ID with brand information
func (s *EventService) GetByIDWithBrand(ctx context.Context, id uuid.UUID) (*EventWithBrand, error) {
	return s.repo.GetByIDWithBrand(ctx, id)
}

// Update updates an existing event
func (s *EventService) Update(ctx context.Context, id uuid.UUID, req *UpdateEventRequest) (*Event, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify event exists
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrStartDateInvalid
	}

	// Update fields
	event.EventType = req.EventType
	event.Organizer = req.Organizer
	event.Name = req.Name
	event.StartDate = startDate
	event.Year = req.Year
	event.DurationDays = req.DurationDays
	event.State = req.State
	event.City = req.City
	event.Venue = req.Venue
	event.Dealer = req.Dealer
	event.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

// ChangeStatus changes the status of an event with validation
func (s *EventService) ChangeStatus(ctx context.Context, id uuid.UUID, req *ChangeStatusRequest) (*Event, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get current event
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newStatus := EventStatus(req.Status)

	// Validate status transition
	if !isValidTransition(event.Status, newStatus) {
		return nil, ErrInvalidStatusTransition
	}

	// Update status
	if err := s.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}

	event.Status = newStatus
	event.UpdatedAt = time.Now()

	return event, nil
}

// Delete deletes an event
func (s *EventService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves events with optional filters
func (s *EventService) List(ctx context.Context, filters EventFilters) ([]*EventWithBrand, error) {
	return s.repo.List(ctx, filters)
}

// ListByBrandID retrieves all events for a specific brand
func (s *EventService) ListByBrandID(ctx context.Context, brandID uuid.UUID) ([]*Event, error) {
	return s.repo.ListByBrandID(ctx, brandID)
}

// isValidTransition validates if a status transition is allowed
func isValidTransition(current, new EventStatus) bool {
	// Define valid transitions
	validTransitions := map[EventStatus][]EventStatus{
		StatusPlanned:   {StatusActive},
		StatusActive:    {StatusCompleted},
		StatusCompleted: {StatusClosed},
		StatusClosed:    {}, // No transitions from CLOSED
	}

	// Allow staying in the same status (no-op)
	if current == new {
		return true
	}

	// Check if transition is allowed
	allowed := validTransitions[current]
	for _, allowedStatus := range allowed {
		if allowedStatus == new {
			return true
		}
	}

	return false
}
