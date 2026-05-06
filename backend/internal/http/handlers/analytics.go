package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/customermx/backend/internal/domain/analytics"
	"github.com/customermx/backend/internal/http/middleware"
)

// AnalyticsHandler maneja las solicitudes HTTP de analytics
type AnalyticsHandler struct {
	analyticsService analytics.Service
}

// NewAnalyticsHandler crea una nueva instancia del handler
func NewAnalyticsHandler(analyticsService analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboard maneja GET /api/v1/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse filters from query params
	filters := analytics.AnalyticsFilters{}

	// Brand ID filter (BRAND users solo ven su marca)
	if claims.Role == "BRAND" && claims.BrandID != nil {
		filters.BrandID = claims.BrandID
	} else if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
		brandID, err := uuid.Parse(brandIDStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid brand_id")
			return
		}
		filters.BrandID = &brandID
	}

	// Year filter
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid year")
			return
		}
		filters.Year = &year
	}

	// Event type filter
	if eventType := r.URL.Query().Get("event_type"); eventType != "" {
		filters.EventType = &eventType
	}

	// Organizer filter
	if organizer := r.URL.Query().Get("organizer"); organizer != "" {
		filters.Organizer = &organizer
	}

	// Setup vendor filter
	if setupVendor := r.URL.Query().Get("setup_vendor"); setupVendor != "" {
		filters.SetupVendor = &setupVendor
	}

	// Get analytics
	result, err := h.analyticsService.GetDashboardAnalytics(ctx, filters)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, http.StatusOK, result)
}

// GetSetupVendors maneja GET /api/v1/analytics/setup-vendors
func (h *AnalyticsHandler) GetSetupVendors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var brandID *uuid.UUID
	if claims.Role == "BRAND" && claims.BrandID != nil {
		brandID = claims.BrandID
	} else if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
		id, err := uuid.Parse(brandIDStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid brand_id")
			return
		}
		brandID = &id
	}

	vendors, err := h.analyticsService.GetAvailableSetupVendors(ctx, brandID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, http.StatusOK, vendors)
}

// GetEventsByBrand maneja GET /api/v1/analytics/events/by-brand
func (h *AnalyticsHandler) GetEventsByBrand(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse filters
	filters := analytics.AnalyticsFilters{}

	// Brand ID filter
	if claims.Role == "BRAND" && claims.BrandID != nil {
		filters.BrandID = claims.BrandID
	} else if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
		brandID, err := uuid.Parse(brandIDStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid brand_id")
			return
		}
		filters.BrandID = &brandID
	}

	// Year filter
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid year")
			return
		}
		filters.Year = &year
	}

	// Get metrics
	result, err := h.analyticsService.GetEventsByBrand(ctx, filters)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, http.StatusOK, result)
}

// GetEventTimeline maneja GET /api/v1/analytics/events/timeline
func (h *AnalyticsHandler) GetEventTimeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse filters
	filters := analytics.AnalyticsFilters{}

	// Brand ID filter
	if claims.Role == "BRAND" && claims.BrandID != nil {
		filters.BrandID = claims.BrandID
	} else if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
		brandID, err := uuid.Parse(brandIDStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid brand_id")
			return
		}
		filters.BrandID = &brandID
	}

	// Year filter
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid year")
			return
		}
		filters.Year = &year
	}

	// Get timeline
	result, err := h.analyticsService.GetEventTimeline(ctx, filters)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, http.StatusOK, result)
}
