package eventphoto

import (
	"time"

	"github.com/google/uuid"
)

const MaxPhotosPerEvent = 10

// EventPhoto represents a photo attached to an event
type EventPhoto struct {
	ID          uuid.UUID `json:"id"`
	EventID     uuid.UUID `json:"event_id"`
	S3Key       string    `json:"-"` // internal, not exposed to clients
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	SizeBytes   int64     `json:"size_bytes"`
	SortOrder   int       `json:"sort_order"`
	UploadedAt  time.Time `json:"uploaded_at"`
	URL         string    `json:"url"` // computed: proxy URL through our API
}
