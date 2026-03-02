-- V7: Create event_photos table
-- Description: Store S3 photo references for events (max 10 per event)

CREATE TABLE event_photos (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id     UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    s3_key       TEXT NOT NULL,
    filename     TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes   BIGINT,
    sort_order   INT DEFAULT 0,
    uploaded_at  TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_event_photos_event ON event_photos(event_id);

COMMENT ON TABLE event_photos IS 'Photos attached to events, stored in S3';
COMMENT ON COLUMN event_photos.s3_key IS 'S3 object key: events/{event_id}/{photo_id}';
COMMENT ON COLUMN event_photos.sort_order IS 'Display order within the event (0-indexed)';
