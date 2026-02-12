-- V3: Create events, coordinators, reports, and notifications tables
-- Description: Event management system with reports and notifications

-- Event status enum
CREATE TYPE event_status AS ENUM ('PLANNED', 'ACTIVE', 'COMPLETED', 'CLOSED');

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id),
    event_type TEXT NOT NULL,
    organizer TEXT NOT NULL,
    name TEXT NOT NULL,
    start_date DATE NOT NULL,
    year INT NOT NULL,
    duration_days INT NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL,
    venue TEXT NOT NULL,
    dealer TEXT NOT NULL,
    status event_status DEFAULT 'PLANNED',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Event coordinators (many-to-many relationship)
CREATE TABLE event_coordinators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT now(),
    UNIQUE(event_id, user_id)
);

-- Event reports (one-to-one with events)
CREATE TABLE event_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
    hostess_count INT,
    setup_vendor TEXT,
    has_promotional BOOLEAN,
    attendees INT,
    activities_count INT,
    leads_collected INT,
    prospects INT,
    dealer_rating INT CHECK (dealer_rating BETWEEN 1 AND 5),
    comments TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Event vehicles (which vehicles were presented at each event)
CREATE TABLE event_vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(event_id, vehicle_id)
);

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    payload JSONB,
    read BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);

-- Indexes for events (critical for performance)
CREATE INDEX idx_events_brand ON events(brand_id);
CREATE INDEX idx_events_year ON events(year);
CREATE INDEX idx_events_date ON events(start_date);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_state_city ON events(state, city);
CREATE INDEX idx_events_brand_year ON events(brand_id, year);
CREATE INDEX idx_events_brand_date ON events(brand_id, start_date);
CREATE INDEX idx_events_brand_status ON events(brand_id, status);

-- Indexes for event coordinators
CREATE INDEX idx_event_coord_event ON event_coordinators(event_id);
CREATE INDEX idx_event_coord_user ON event_coordinators(user_id);

-- Indexes for event reports
CREATE INDEX idx_event_reports_event ON event_reports(event_id);
CREATE INDEX idx_event_reports_rating ON event_reports(dealer_rating);
CREATE INDEX idx_event_reports_completed ON event_reports(completed);

-- Indexes for event vehicles
CREATE INDEX idx_event_vehicles_event ON event_vehicles(event_id);
CREATE INDEX idx_event_vehicles_vehicle ON event_vehicles(vehicle_id);

-- Indexes for notifications
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(read);
CREATE INDEX idx_notifications_type ON notifications(type);

-- Comments for documentation
COMMENT ON TABLE events IS 'Promotional events for automotive brands';
COMMENT ON COLUMN events.event_type IS 'Type of event (Triatlón, Fútbol, Golf, etc.)';
COMMENT ON COLUMN events.status IS 'Event lifecycle: PLANNED -> ACTIVE -> COMPLETED -> CLOSED';
COMMENT ON TABLE event_coordinators IS 'Coordinators assigned to each event';
COMMENT ON TABLE event_reports IS 'Operational data collected during and after the event';
COMMENT ON TABLE event_vehicles IS 'Vehicles presented at each event with quantities';
COMMENT ON COLUMN event_vehicles.quantity IS 'Number of units of this vehicle model at the event';
COMMENT ON TABLE notifications IS 'User notifications for event assignments and completions';
COMMENT ON COLUMN notifications.type IS 'Notification type (event_assigned, event_completed, etc.)';
