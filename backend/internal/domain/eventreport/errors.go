package eventreport

import "errors"

// Validation errors
var (
	ErrDealerRatingInvalid = errors.New("dealer rating must be between 1 and 5")
)

// Business logic errors
var (
	ErrReportNotFound    = errors.New("event report not found")
	ErrEventNotFound     = errors.New("event not found")
	ErrEventNotCompleted = errors.New("event must be in COMPLETED or CLOSED status to mark report as completed")
)
