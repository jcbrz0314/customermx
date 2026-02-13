package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/customermx/backend/internal/domain/event"
	"github.com/customermx/backend/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EventHandler handles event endpoints
type EventHandler struct {
	eventService event.Service
}

// NewEventHandler creates a new EventHandler
func NewEventHandler(eventService event.Service) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEvent creates a new event
// POST /api/v1/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req event.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eventResp, err := h.eventService.Create(r.Context(), &req)
	if err != nil {
		switch err {
		case event.ErrBrandIDRequired, event.ErrNameRequired, event.ErrEventTypeRequired,
			event.ErrStartDateRequired, event.ErrStartDateInvalid, event.ErrYearInvalid,
			event.ErrDurationInvalid:
			RespondError(w, http.StatusBadRequest, err.Error())
		case event.ErrBrandNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create event")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, eventResp)
}

// ListEvents retrieves events with optional filters
// GET /api/v1/events?brand_id=x&year=y&status=z&state=s
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	filters := event.EventFilters{}

	// Get user claims for role-based filtering
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// BRAND users: automatically filter to their brand
	if claims.Role == "BRAND" && claims.BrandID != nil {
		filters.BrandID = claims.BrandID
	}

	// Parse brand_id filter from query (ADMIN and COORDINATOR can filter by brand)
	if claims.Role != "BRAND" {
		if brandIDStr := r.URL.Query().Get("brand_id"); brandIDStr != "" {
			brandID, err := uuid.Parse(brandIDStr)
			if err != nil {
				RespondError(w, http.StatusBadRequest, "Invalid brand_id")
				return
			}
			filters.BrandID = &brandID
		}
	}

	// Parse year filter
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid year")
			return
		}
		filters.Year = &year
	}

	// Parse status filter
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := event.EventStatus(statusStr)
		filters.Status = &status
	}

	// Parse state filter
	if stateStr := r.URL.Query().Get("state"); stateStr != "" {
		filters.State = &stateStr
	}

	events, err := h.eventService.List(r.Context(), filters)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}

	RespondSuccess(w, http.StatusOK, events)
}

// GetEvent retrieves an event by ID
// GET /api/v1/events/{id}
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	eventResp, err := h.eventService.GetByIDWithBrand(r.Context(), id)
	if err != nil {
		if err == event.ErrEventNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get event")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, eventResp)
}

// UpdateEvent updates an existing event
// PUT /api/v1/events/{id}
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req event.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eventResp, err := h.eventService.Update(r.Context(), id, &req)
	if err != nil {
		switch err {
		case event.ErrEventNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case event.ErrNameRequired, event.ErrEventTypeRequired, event.ErrStartDateRequired,
			event.ErrStartDateInvalid, event.ErrYearInvalid, event.ErrDurationInvalid:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to update event")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, eventResp)
}

// ChangeEventStatus changes the status of an event
// PATCH /api/v1/events/{id}/status
func (h *EventHandler) ChangeEventStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req event.ChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eventResp, err := h.eventService.ChangeStatus(r.Context(), id, &req)
	if err != nil {
		switch err {
		case event.ErrEventNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case event.ErrInvalidStatus, event.ErrInvalidStatusTransition:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to change event status")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, eventResp)
}

// DeleteEvent deletes an event
// DELETE /api/v1/events/{id}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := h.eventService.Delete(r.Context(), id); err != nil {
		if err == event.ErrEventNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete event")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Event deleted successfully")
}

// ListEventsByBrand retrieves all events for a specific brand
// GET /api/v1/brands/{brandId}/events
func (h *EventHandler) ListEventsByBrand(w http.ResponseWriter, r *http.Request) {
	brandIDStr := chi.URLParam(r, "brandId")
	brandID, err := uuid.Parse(brandIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	events, err := h.eventService.ListByBrandID(r.Context(), brandID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}

	RespondSuccess(w, http.StatusOK, events)
}
