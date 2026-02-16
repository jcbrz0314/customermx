package analytics

import (
	"context"
)

// Service define las operaciones de negocio para analytics
type Service interface {
	GetDashboardAnalytics(ctx context.Context, filters AnalyticsFilters) (*DashboardAnalytics, error)
	GetEventsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error)
	GetEventTimeline(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error)
}

type service struct {
	repo Repository
}

// NewService crea una nueva instancia del servicio
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetDashboardAnalytics obtiene todas las métricas del dashboard
func (s *service) GetDashboardAnalytics(ctx context.Context, filters AnalyticsFilters) (*DashboardAnalytics, error) {
	// Obtener todas las métricas en paralelo usando channels
	type result struct {
		totals         *TotalMetrics
		byBrand        []BrandMetrics
		byMonth        []MonthlyMetrics
		byState        []StateMetrics
		topVehicles    []VehicleMetrics
		yearComparison []YearComparison
		byEventType    []EventTypeMetrics
		topDealers     []DealerMetrics
		conversion     *ConversionMetrics
		topCities      []CityMetrics
		err            error
	}

	ch := make(chan result, 1)

	go func() {
		var r result

		// Total metrics
		r.totals, r.err = s.repo.GetTotalMetrics(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Metrics by brand
		r.byBrand, r.err = s.repo.GetMetricsByBrand(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Metrics by month
		r.byMonth, r.err = s.repo.GetMetricsByMonth(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Metrics by state
		r.byState, r.err = s.repo.GetMetricsByState(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Top vehicles (limit 10)
		r.topVehicles, r.err = s.repo.GetTopVehicles(ctx, filters, 10)
		if r.err != nil {
			ch <- r
			return
		}

		// Year comparison
		r.yearComparison, r.err = s.repo.GetYearComparison(ctx, filters.BrandID)
		if r.err != nil {
			ch <- r
			return
		}

		// Metrics by event type
		r.byEventType, r.err = s.repo.GetMetricsByEventType(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Top dealers
		r.topDealers, r.err = s.repo.GetMetricsByDealer(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Conversion metrics
		r.conversion, r.err = s.repo.GetConversionMetrics(ctx, filters)
		if r.err != nil {
			ch <- r
			return
		}

		// Top cities (limit 10)
		r.topCities, r.err = s.repo.GetMetricsByCity(ctx, filters, 10)
		if r.err != nil {
			ch <- r
			return
		}

		ch <- r
	}()

	r := <-ch
	if r.err != nil {
		return nil, r.err
	}

	analytics := &DashboardAnalytics{
		Totals:         *r.totals,
		ByBrand:        r.byBrand,
		ByMonth:        r.byMonth,
		ByState:        r.byState,
		TopVehicles:    r.topVehicles,
		YearComparison: r.yearComparison,
		ByEventType:    r.byEventType,
		TopDealers:     r.topDealers,
		Conversion:     r.conversion,
		TopCities:      r.topCities,
	}

	// Inicializar slices vacíos si son nil
	if analytics.ByBrand == nil {
		analytics.ByBrand = []BrandMetrics{}
	}
	if analytics.ByMonth == nil {
		analytics.ByMonth = []MonthlyMetrics{}
	}
	if analytics.ByState == nil {
		analytics.ByState = []StateMetrics{}
	}
	if analytics.TopVehicles == nil {
		analytics.TopVehicles = []VehicleMetrics{}
	}
	if analytics.YearComparison == nil {
		analytics.YearComparison = []YearComparison{}
	}
	if analytics.ByEventType == nil {
		analytics.ByEventType = []EventTypeMetrics{}
	}
	if analytics.TopDealers == nil {
		analytics.TopDealers = []DealerMetrics{}
	}
	if analytics.TopCities == nil {
		analytics.TopCities = []CityMetrics{}
	}

	return analytics, nil
}

// GetEventsByBrand obtiene métricas agrupadas por marca
func (s *service) GetEventsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error) {
	return s.repo.GetMetricsByBrand(ctx, filters)
}

// GetEventTimeline obtiene timeline de eventos por mes
func (s *service) GetEventTimeline(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error) {
	return s.repo.GetMetricsByMonth(ctx, filters)
}
