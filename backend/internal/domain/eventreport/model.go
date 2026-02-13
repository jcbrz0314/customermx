package eventreport

import (
	"time"

	"github.com/google/uuid"
)

// EventReport represents operational metrics from an event
type EventReport struct {
	ID               uuid.UUID  `json:"id"`
	EventID          uuid.UUID  `json:"event_id"`
	HostessCount     *int       `json:"hostess_count,omitempty"`
	SetupVendor      *string    `json:"setup_vendor,omitempty"`
	HasPromotional   *bool      `json:"has_promotional,omitempty"`
	Attendees        *int       `json:"attendees,omitempty"`
	ActivitiesCount  *int       `json:"activities_count,omitempty"`
	LeadsCollected   *int       `json:"leads_collected,omitempty"`
	Prospects        *int       `json:"prospects,omitempty"`
	DealerRating     *int       `json:"dealer_rating,omitempty"`     // 1-5 scale
	Comments         *string    `json:"comments,omitempty"`
	Completed        bool       `json:"completed"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CreateOrUpdateReportRequest represents the request to create or update a report
type CreateOrUpdateReportRequest struct {
	HostessCount    *int    `json:"hostess_count,omitempty"`
	SetupVendor     *string `json:"setup_vendor,omitempty"`
	HasPromotional  *bool   `json:"has_promotional,omitempty"`
	Attendees       *int    `json:"attendees,omitempty"`
	ActivitiesCount *int    `json:"activities_count,omitempty"`
	LeadsCollected  *int    `json:"leads_collected,omitempty"`
	Prospects       *int    `json:"prospects,omitempty"`
	DealerRating    *int    `json:"dealer_rating,omitempty"`
	Comments        *string `json:"comments,omitempty"`
}

// Validate validates the create or update report request
func (r *CreateOrUpdateReportRequest) Validate() error {
	// Validate dealer rating if provided
	if r.DealerRating != nil {
		if *r.DealerRating < 1 || *r.DealerRating > 5 {
			return ErrDealerRatingInvalid
		}
	}
	return nil
}

// CompleteReportRequest represents the request to mark a report as completed
type CompleteReportRequest struct {
	Completed bool `json:"completed"`
}
