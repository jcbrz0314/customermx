package analytics

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository define las operaciones de analytics
type Repository interface {
	GetTotalMetrics(ctx context.Context, filters AnalyticsFilters) (*TotalMetrics, error)
	GetMetricsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error)
	GetMetricsByMonth(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error)
	GetMetricsByState(ctx context.Context, filters AnalyticsFilters) ([]StateMetrics, error)
	GetTopVehicles(ctx context.Context, filters AnalyticsFilters, limit int) ([]VehicleMetrics, error)
	GetYearComparison(ctx context.Context, brandID *uuid.UUID) ([]YearComparison, error)
	GetMetricsByEventType(ctx context.Context, filters AnalyticsFilters) ([]EventTypeMetrics, error)
	GetMetricsByDealer(ctx context.Context, filters AnalyticsFilters) ([]DealerMetrics, error)
	GetConversionMetrics(ctx context.Context, filters AnalyticsFilters) (*ConversionMetrics, error)
	GetMetricsByCity(ctx context.Context, filters AnalyticsFilters, limit int) ([]CityMetrics, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

// GetTotalMetrics obtiene métricas totales agregadas
func (r *postgresRepository) GetTotalMetrics(ctx context.Context, filters AnalyticsFilters) (*TotalMetrics, error) {
	query := `
		SELECT
			COUNT(e.id) as total_events,
			COALESCE(SUM(er.attendees), 0) as total_attendees,
			COALESCE(SUM(er.leads_collected), 0) as total_leads,
			COALESCE(SUM(er.prospects), 0) as total_prospects,
			COALESCE(AVG(er.attendees), 0) as average_attendees,
			COALESCE(AVG(er.dealer_rating), 0) as average_rating
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
	`

	var metrics TotalMetrics
	err := r.pool.QueryRow(ctx, query, filters.BrandID, filters.Year).Scan(
		&metrics.TotalEvents,
		&metrics.TotalAttendees,
		&metrics.TotalLeads,
		&metrics.TotalProspects,
		&metrics.AverageAttendees,
		&metrics.AverageRating,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting total metrics: %w", err)
	}

	return &metrics, nil
}

// GetMetricsByBrand obtiene métricas agrupadas por marca
func (r *postgresRepository) GetMetricsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error) {
	query := `
		SELECT
			b.id as brand_id,
			b.name as brand_name,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as total_attendees,
			COALESCE(SUM(er.leads_collected), 0) as total_leads,
			COALESCE(AVG(er.dealer_rating), 0) as average_rating
		FROM brands b
		LEFT JOIN events e ON b.id = e.brand_id AND e.status = 'COMPLETED'
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR b.id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		GROUP BY b.id, b.name
		HAVING COUNT(e.id) > 0
		ORDER BY event_count DESC
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by brand: %w", err)
	}
	defer rows.Close()

	var metrics []BrandMetrics
	for rows.Next() {
		var m BrandMetrics
		err := rows.Scan(
			&m.BrandID,
			&m.BrandName,
			&m.EventCount,
			&m.TotalAttendees,
			&m.TotalLeads,
			&m.AverageRating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning brand metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating brand metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetMetricsByMonth obtiene métricas mensuales para timeline
func (r *postgresRepository) GetMetricsByMonth(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error) {
	query := `
		SELECT
			EXTRACT(YEAR FROM e.start_date)::int as year,
			EXTRACT(MONTH FROM e.start_date)::int as month,
			TO_CHAR(e.start_date, 'Month') as month_name,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as attendees
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
		GROUP BY year, month, month_name
		ORDER BY year, month
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by month: %w", err)
	}
	defer rows.Close()

	var metrics []MonthlyMetrics
	for rows.Next() {
		var m MonthlyMetrics
		err := rows.Scan(
			&m.Year,
			&m.Month,
			&m.MonthName,
			&m.EventCount,
			&m.Attendees,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning monthly metrics: %w", err)
		}
		// Trim whitespace from month name
		m.MonthName = trimSpace(m.MonthName)
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating monthly metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetMetricsByState obtiene métricas por estado geográfico
func (r *postgresRepository) GetMetricsByState(ctx context.Context, filters AnalyticsFilters) ([]StateMetrics, error) {
	query := `
		SELECT
			e.state,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as attendees
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
		GROUP BY e.state
		ORDER BY event_count DESC
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by state: %w", err)
	}
	defer rows.Close()

	var metrics []StateMetrics
	for rows.Next() {
		var m StateMetrics
		err := rows.Scan(
			&m.State,
			&m.EventCount,
			&m.Attendees,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning state metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating state metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetTopVehicles obtiene los vehículos más presentados
func (r *postgresRepository) GetTopVehicles(ctx context.Context, filters AnalyticsFilters, limit int) ([]VehicleMetrics, error) {
	query := `
		SELECT
			v.id as vehicle_id,
			v.model_name,
			b.name as brand_name,
			COUNT(DISTINCT ev.event_id) as times_presented,
			SUM(ev.quantity) as total_quantity
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.id
		JOIN event_vehicles ev ON v.id = ev.vehicle_id
		JOIN events e ON ev.event_id = e.id
		WHERE e.status = 'COMPLETED'
		  AND ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		GROUP BY v.id, v.model_name, b.name
		ORDER BY total_quantity DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting top vehicles: %w", err)
	}
	defer rows.Close()

	var metrics []VehicleMetrics
	for rows.Next() {
		var m VehicleMetrics
		err := rows.Scan(
			&m.VehicleID,
			&m.ModelName,
			&m.BrandName,
			&m.TimesPresented,
			&m.TotalQuantity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating vehicle metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetYearComparison obtiene comparativa año vs año
func (r *postgresRepository) GetYearComparison(ctx context.Context, brandID *uuid.UUID) ([]YearComparison, error) {
	query := `
		SELECT
			e.year,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as total_attendees,
			COALESCE(AVG(er.dealer_rating), 0) as average_rating
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND e.status = 'COMPLETED'
		GROUP BY e.year
		ORDER BY e.year
	`

	rows, err := r.pool.Query(ctx, query, brandID)
	if err != nil {
		return nil, fmt.Errorf("error getting year comparison: %w", err)
	}
	defer rows.Close()

	var comparisons []YearComparison
	for rows.Next() {
		var c YearComparison
		err := rows.Scan(
			&c.Year,
			&c.EventCount,
			&c.TotalAttendees,
			&c.AverageRating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning year comparison: %w", err)
		}
		comparisons = append(comparisons, c)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating year comparison: %w", rows.Err())
	}

	return comparisons, nil
}

// GetMetricsByEventType obtiene métricas por tipo de evento
func (r *postgresRepository) GetMetricsByEventType(ctx context.Context, filters AnalyticsFilters) ([]EventTypeMetrics, error) {
	query := `
		SELECT
			e.event_type,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as attendees
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
		GROUP BY e.event_type
		ORDER BY event_count DESC
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by event type: %w", err)
	}
	defer rows.Close()

	var metrics []EventTypeMetrics
	for rows.Next() {
		var m EventTypeMetrics
		err := rows.Scan(
			&m.EventType,
			&m.EventCount,
			&m.Attendees,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event type metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating event type metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetMetricsByDealer obtiene ranking de distribuidores
func (r *postgresRepository) GetMetricsByDealer(ctx context.Context, filters AnalyticsFilters) ([]DealerMetrics, error) {
	query := `
		SELECT
			e.dealer,
			COUNT(e.id) as event_count,
			COALESCE(AVG(er.dealer_rating), 0) as average_rating,
			COALESCE(SUM(er.attendees), 0) as total_attendees
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
		  AND e.dealer IS NOT NULL AND e.dealer != ''
		GROUP BY e.dealer
		ORDER BY average_rating DESC, event_count DESC
		LIMIT 10
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by dealer: %w", err)
	}
	defer rows.Close()

	var metrics []DealerMetrics
	for rows.Next() {
		var m DealerMetrics
		err := rows.Scan(
			&m.Dealer,
			&m.EventCount,
			&m.AverageRating,
			&m.TotalAttendees,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning dealer metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating dealer metrics: %w", rows.Err())
	}

	return metrics, nil
}

// GetConversionMetrics obtiene tasas de conversión
func (r *postgresRepository) GetConversionMetrics(ctx context.Context, filters AnalyticsFilters) (*ConversionMetrics, error) {
	query := `
		SELECT
			COALESCE(SUM(er.attendees), 0) as total_attendees,
			COALESCE(SUM(er.leads_collected), 0) as total_leads,
			COALESCE(SUM(er.prospects), 0) as total_prospects,
			CASE
				WHEN SUM(er.attendees) > 0 THEN
					(SUM(er.leads_collected)::float / SUM(er.attendees)::float) * 100
				ELSE 0
			END as lead_rate,
			CASE
				WHEN SUM(er.leads_collected) > 0 THEN
					(SUM(er.prospects)::float / SUM(er.leads_collected)::float) * 100
				ELSE 0
			END as prospect_rate
		FROM events e
		JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
	`

	var metrics ConversionMetrics
	err := r.pool.QueryRow(ctx, query, filters.BrandID, filters.Year).Scan(
		&metrics.TotalAttendees,
		&metrics.TotalLeads,
		&metrics.TotalProspects,
		&metrics.LeadRate,
		&metrics.ProspectRate,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting conversion metrics: %w", err)
	}

	return &metrics, nil
}

// GetMetricsByCity obtiene métricas por ciudad (top N)
func (r *postgresRepository) GetMetricsByCity(ctx context.Context, filters AnalyticsFilters, limit int) ([]CityMetrics, error) {
	query := `
		SELECT
			e.state,
			e.city,
			COUNT(e.id) as event_count,
			COALESCE(SUM(er.attendees), 0) as attendees
		FROM events e
		LEFT JOIN event_reports er ON e.id = er.event_id
		WHERE ($1::uuid IS NULL OR e.brand_id = $1)
		  AND ($2::int IS NULL OR e.year = $2)
		  AND e.status = 'COMPLETED'
		GROUP BY e.state, e.city
		ORDER BY event_count DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, filters.BrandID, filters.Year, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting metrics by city: %w", err)
	}
	defer rows.Close()

	var metrics []CityMetrics
	for rows.Next() {
		var m CityMetrics
		err := rows.Scan(
			&m.State,
			&m.City,
			&m.EventCount,
			&m.Attendees,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning city metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating city metrics: %w", rows.Err())
	}

	return metrics, nil
}

// trimSpace es un helper para limpiar espacios en blanco
func trimSpace(s string) string {
	result := ""
	for _, c := range s {
		if c != ' ' {
			result += string(c)
		}
	}
	return result
}
