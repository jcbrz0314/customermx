package eventcoordinator

import (
	"context"
	"fmt"
	"time"

	"github.com/customermx/backend/internal/domain/event"
	"github.com/customermx/backend/internal/domain/user"
	"github.com/customermx/backend/internal/infra/mail"
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
	repo        Repository
	eventRepo   event.Repository
	userRepo    user.Repository
	mailService mail.Service
}

// NewService creates a new EventCoordinatorService
func NewService(repo Repository, eventRepo event.Repository, userRepo user.Repository, mailService mail.Service) *EventCoordinatorService {
	return &EventCoordinatorService{
		repo:        repo,
		eventRepo:   eventRepo,
		userRepo:    userRepo,
		mailService: mailService,
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

	// Fetch event details while the request context is still alive
	evt, err := s.eventRepo.GetByIDWithBrand(ctx, eventID)
	if err == nil {
		details := mail.EventAssignmentDetails{
			EventID:      evt.ID.String(),
			EventName:    evt.Name,
			BrandName:    evt.BrandName,
			EventType:    evt.EventType,
			Organizer:    evt.Organizer,
			StartDate:    evt.StartDate.Format("02/01/2006"),
			DurationDays: evt.DurationDays,
			State:        evt.State,
			City:         evt.City,
			Venue:        evt.Venue,
			Dealer:       evt.Dealer,
		}
		recipientEmail := usr.Email
		// Use context.Background() so the goroutine outlives the HTTP request
		go func() {
			if err := s.mailService.SendEventAssignment(recipientEmail, details); err != nil {
				fmt.Printf("[mail] error sending assignment email to %s: %v\n", recipientEmail, err)
			}
		}()
	}

	return ec, nil
}

// Remove removes a coordinator from an event and notifies them by email
func (s *EventCoordinatorService) Remove(ctx context.Context, eventID, userID uuid.UUID) error {
	// Fetch user and event details before removing (while context is still valid)
	usr, userErr := s.userRepo.GetByID(ctx, userID)
	evt, evtErr := s.eventRepo.GetByIDWithBrand(ctx, eventID)

	if err := s.repo.Remove(ctx, eventID, userID); err != nil {
		return err
	}

	// Send unassignment email if we have all the data
	if userErr == nil && evtErr == nil {
		details := mail.EventUnassignmentDetails{
			EventName: evt.Name,
			BrandName: evt.BrandName,
			StartDate: evt.StartDate.Format("02/01/2006"),
			City:      evt.City,
			State:     evt.State,
		}
		recipientEmail := usr.Email
		go func() {
			if err := s.mailService.SendEventUnassignment(recipientEmail, details); err != nil {
				fmt.Printf("[mail] error sending unassignment email to %s: %v\n", recipientEmail, err)
			}
		}()
	}

	return nil
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
