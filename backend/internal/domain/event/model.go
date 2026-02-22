package event

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventStatus represents the status of an event
type EventStatus string

const (
	StatusPlanned   EventStatus = "PLANNED"
	StatusActive    EventStatus = "ACTIVE"
	StatusCompleted EventStatus = "COMPLETED"
	StatusClosed    EventStatus = "CLOSED"
)

// Event represents an automotive promotional event
type Event struct {
	ID           uuid.UUID   `json:"id"`
	BrandID      uuid.UUID   `json:"brand_id"`
	EventType    string      `json:"event_type"`
	Organizer    string      `json:"organizer"`
	Name         string      `json:"name"`
	StartDate    time.Time   `json:"start_date"`
	Year         int         `json:"year"`
	DurationDays int         `json:"duration_days"`
	State        string      `json:"state"`
	City         string      `json:"city"`
	Venue        string      `json:"venue"`
	Dealer       string      `json:"dealer"`
	Status       EventStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// EventWithBrand includes the brand name via JOIN
type EventWithBrand struct {
	Event
	BrandName string `json:"brand_name"`
}

// CreateEventRequest represents the request to create an event
type CreateEventRequest struct {
	BrandID      string `json:"brand_id"`
	EventType    string `json:"event_type"`
	Organizer    string `json:"organizer"`
	Name         string `json:"name"`
	StartDate    string `json:"start_date"` // Format: YYYY-MM-DD
	Year         int    `json:"year"`
	DurationDays int    `json:"duration_days"`
	State        string `json:"state"`
	City         string `json:"city"`
	Venue        string `json:"venue"`
	Dealer       string `json:"dealer"`
}

// Validate validates the create event request
func (r *CreateEventRequest) Validate() error {
	if r.BrandID == "" {
		return ErrBrandIDRequired
	}
	if _, err := uuid.Parse(r.BrandID); err != nil {
		return ErrBrandIDRequired
	}
	if strings.TrimSpace(r.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(r.EventType) == "" {
		return ErrEventTypeRequired
	}
	if r.StartDate == "" {
		return ErrStartDateRequired
	}
	if _, err := time.Parse("2006-01-02", r.StartDate); err != nil {
		return ErrStartDateInvalid
	}
	if r.Year < 2000 || r.Year > 2100 {
		return ErrYearInvalid
	}
	if r.DurationDays < 1 {
		return ErrDurationInvalid
	}
	return nil
}

// UpdateEventRequest represents the request to update an event
type UpdateEventRequest struct {
	EventType    string `json:"event_type"`
	Organizer    string `json:"organizer"`
	Name         string `json:"name"`
	StartDate    string `json:"start_date"` // Format: YYYY-MM-DD
	Year         int    `json:"year"`
	DurationDays int    `json:"duration_days"`
	State        string `json:"state"`
	City         string `json:"city"`
	Venue        string `json:"venue"`
	Dealer       string `json:"dealer"`
}

// Validate validates the update event request
func (r *UpdateEventRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(r.EventType) == "" {
		return ErrEventTypeRequired
	}
	if r.StartDate == "" {
		return ErrStartDateRequired
	}
	if _, err := time.Parse("2006-01-02", r.StartDate); err != nil {
		return ErrStartDateInvalid
	}
	if r.Year < 2000 || r.Year > 2100 {
		return ErrYearInvalid
	}
	if r.DurationDays < 1 {
		return ErrDurationInvalid
	}
	return nil
}

// ChangeStatusRequest represents the request to change event status
type ChangeStatusRequest struct {
	Status string `json:"status"`
}

// Validate validates the change status request
func (r *ChangeStatusRequest) Validate() error {
	status := EventStatus(r.Status)
	if status != StatusPlanned && status != StatusActive &&
		status != StatusCompleted && status != StatusClosed {
		return ErrInvalidStatus
	}
	return nil
}

// EventFilters represents filters for querying events
type EventFilters struct {
	BrandID       *uuid.UUID
	Year          *int
	Status        *EventStatus
	State         *string
	CoordinatorID *uuid.UUID
}
