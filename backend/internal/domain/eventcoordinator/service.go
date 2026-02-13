package eventcoordinator

import (
	"context"
	"time"

	"github.com/customermx/backend/internal/domain/event"
	"github.com/customermx/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Service defines the interface for event coordinator business logic
type Service interface {
	Assign(ctx context.Context, req *AssignCoordinatorRequest) (*EventCoordinator, error)
	Remove(ctx context.Context, eventID, userID uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventCoordinatorWithUser, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*EventCoordinator, error)
	IsCoordinatorAssigned(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}

// EventCoordinatorService implements the Service interface
type EventCoordinatorService struct {
	repo      Repository
	eventRepo event.Repository
	userRepo  user.Repository
}

// NewService creates a new EventCoordinatorService
func NewService(repo Repository, eventRepo event.Repository, userRepo user.Repository) *EventCoordinatorService {
	return &EventCoordinatorService{
		repo:      repo,
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

// Assign assigns a coordinator to an event
func (s *EventCoordinatorService) Assign(ctx context.Context, req *AssignCoordinatorRequest) (*EventCoordinator, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return nil, ErrEventIDRequired
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ErrUserIDRequired
	}

	// Verify event exists
	_, err = s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if err == event.ErrEventNotFound {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	// Verify user exists and is a coordinator
	usr, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if usr.Role != user.RoleCoordinator {
		return nil, ErrUserNotCoordinator
	}

	// Create assignment
	ec := &EventCoordinator{
		ID:         uuid.New(),
		EventID:    eventID,
		UserID:     userID,
		AssignedAt: time.Now(),
	}

	if err := s.repo.Assign(ctx, ec); err != nil {
		return nil, err
	}

	return ec, nil
}

// Remove removes a coordinator from an event
func (s *EventCoordinatorService) Remove(ctx context.Context, eventID, userID uuid.UUID) error {
	return s.repo.Remove(ctx, eventID, userID)
}

// ListByEvent retrieves all coordinators for an event
func (s *EventCoordinatorService) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]*EventCoordinatorWithUser, error) {
	return s.repo.ListByEvent(ctx, eventID)
}

// ListByUser retrieves all events assigned to a coordinator
func (s *EventCoordinatorService) ListByUser(ctx context.Context, userID uuid.UUID) ([]*EventCoordinator, error) {
	return s.repo.ListByUser(ctx, userID)
}

// IsCoordinatorAssigned checks if a user is assigned as coordinator to an event
func (s *EventCoordinatorService) IsCoordinatorAssigned(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	return s.repo.IsCoordinatorAssigned(ctx, eventID, userID)
}
