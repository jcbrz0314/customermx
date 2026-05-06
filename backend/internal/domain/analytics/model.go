package analytics

import (
	"time"

	"github.com/google/uuid"
)

// TotalMetrics representa métricas generales agregadas
type TotalMetrics struct {
	TotalEvents      int     `json:"total_events"`
	TotalAttendees   int     `json:"total_attendees"`
	TotalLeads       int     `json:"total_leads"`
	TotalProspects   int     `json:"total_prospects"`
	AverageAttendees float64 `json:"average_attendees"`
	AverageRating    float64 `json:"average_rating"`
}

// BrandMetrics representa métricas agrupadas por marca
type BrandMetrics struct {
	BrandID        uuid.UUID `json:"brand_id"`
	BrandName      string    `json:"brand_name"`
	EventCount     int       `json:"event_count"`
	TotalAttendees int       `json:"total_attendees"`
	TotalLeads     int       `json:"total_leads"`
	AverageRating  float64   `json:"average_rating"`
}

// MonthlyMetrics representa métricas mensuales para timeline
type MonthlyMetrics struct {
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	MonthName  string `json:"month_name"`
	EventCount int    `json:"event_count"`
	Attendees  int    `json:"attendees"`
}

// StateMetrics representa métricas por estado geográfico
type StateMetrics struct {
	State      string `json:"state"`
	EventCount int    `json:"event_count"`
	Attendees  int    `json:"attendees"`
}

// VehicleMetrics representa métricas de vehículos más presentados
type VehicleMetrics struct {
	VehicleID      uuid.UUID `json:"vehicle_id"`
	ModelName      string    `json:"model_name"`
	BrandName      string    `json:"brand_name"`
	TimesPresented int       `json:"times_presented"`
	TotalQuantity  int       `json:"total_quantity"`
}

// YearComparison representa comparativa año vs año
type YearComparison struct {
	Year           int     `json:"year"`
	EventCount     int     `json:"event_count"`
	TotalAttendees int     `json:"total_attendees"`
	AverageRating  float64 `json:"average_rating"`
}

// AnalyticsFilters contiene los filtros opcionales para queries
type AnalyticsFilters struct {
	BrandID     *uuid.UUID `json:"brand_id"`
	Year        *int       `json:"year"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	EventType   *string    `json:"event_type"`
	Organizer   *string    `json:"organizer"`
	SetupVendor *string    `json:"setup_vendor"`
}

// EventTypeMetrics representa métricas por tipo de evento
type EventTypeMetrics struct {
	EventType  string `json:"event_type"`
	EventCount int    `json:"event_count"`
	Attendees  int    `json:"attendees"`
}

// DealerMetrics representa ranking de distribuidores
type DealerMetrics struct {
	Dealer         string  `json:"dealer"`
	EventCount     int     `json:"event_count"`
	AverageRating  float64 `json:"average_rating"`
	TotalAttendees int     `json:"total_attendees"`
}

// ConversionMetrics representa tasas de conversión
type ConversionMetrics struct {
	TotalAttendees int     `json:"total_attendees"`
	TotalLeads     int     `json:"total_leads"`
	TotalProspects int     `json:"total_prospects"`
	LeadRate       float64 `json:"lead_rate"`     // leads / attendees * 100
	ProspectRate   float64 `json:"prospect_rate"` // prospects / leads * 100
}

// CityMetrics representa métricas por ciudad
type CityMetrics struct {
	State      string `json:"state"`
	City       string `json:"city"`
	EventCount int    `json:"event_count"`
	Attendees  int    `json:"attendees"`
}

// VenueMetrics representa métricas agrupadas por sede
type VenueMetrics struct {
	Venue          string `json:"venue"`
	EventCount     int    `json:"event_count"`
	TotalAttendees int    `json:"total_attendees"`
	TotalLeads     int    `json:"total_leads"`
	TotalProspects int    `json:"total_prospects"`
}

// DashboardAnalytics es la respuesta completa del dashboard
type DashboardAnalytics struct {
	Totals         TotalMetrics        `json:"totals"`
	ByBrand        []BrandMetrics      `json:"by_brand"`
	ByMonth        []MonthlyMetrics    `json:"by_month"`
	ByState        []StateMetrics      `json:"by_state"`
	TopVehicles    []VehicleMetrics    `json:"top_vehicles"`
	YearComparison []YearComparison    `json:"year_comparison"`
	ByEventType    []EventTypeMetrics  `json:"by_event_type"`
	TopDealers     []DealerMetrics     `json:"top_dealers"`
	Conversion     *ConversionMetrics  `json:"conversion"`
	TopCities      []CityMetrics       `json:"top_cities"`
	ByVenue        []VenueMetrics      `json:"by_venue"`
}
