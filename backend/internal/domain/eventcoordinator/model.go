package eventcoordinator

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventCoordinator represents the assignment of a coordinator to an event
type EventCoordinator struct {
	ID         uuid.UUID `json:"id"`
	EventID    uuid.UUID `json:"event_id"`
	UserID     uuid.UUID `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// EventCoordinatorWithUser includes user information via JOIN
type EventCoordinatorWithUser struct {
	EventCoordinator
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

// AssignCoordinatorRequest represents the request to assign a coordinator
type AssignCoordinatorRequest struct {
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
}

// Validate validates the assign coordinator request
func (r *AssignCoordinatorRequest) Validate() error {
	if strings.TrimSpace(r.EventID) == "" {
		return ErrEventIDRequired
	}
	if _, err := uuid.Parse(r.EventID); err != nil {
		return ErrEventIDRequired
	}
	if strings.TrimSpace(r.UserID) == "" {
		return ErrUserIDRequired
	}
	if _, err := uuid.Parse(r.UserID); err != nil {
		return ErrUserIDRequired
	}
	return nil
}
