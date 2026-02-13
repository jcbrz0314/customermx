package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/eventvehicle"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EventVehicleHandler handles event vehicle endpoints
type EventVehicleHandler struct {
	vehicleService eventvehicle.Service
}

// NewEventVehicleHandler creates a new EventVehicleHandler
func NewEventVehicleHandler(vehicleService eventvehicle.Service) *EventVehicleHandler {
	return &EventVehicleHandler{vehicleService: vehicleService}
}

// AddVehicle adds a vehicle to an event
// POST /api/v1/events/{eventId}/vehicles
func (h *EventVehicleHandler) AddVehicle(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	var req eventvehicle.AddVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Override eventId from URL if present
	if eventID != "" {
		req.EventID = eventID
	}

	ev, err := h.vehicleService.AddVehicle(r.Context(), &req)
	if err != nil {
		switch err {
		case eventvehicle.ErrEventIDRequired, eventvehicle.ErrVehicleIDRequired,
			eventvehicle.ErrQuantityInvalid:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventvehicle.ErrEventNotFound, eventvehicle.ErrVehicleNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case eventvehicle.ErrBrandMismatch:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventvehicle.ErrVehicleAlreadyAdded:
			RespondError(w, http.StatusConflict, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to add vehicle to event")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, ev)
}

// RemoveVehicle removes a vehicle from an event
// DELETE /api/v1/events/{eventId}/vehicles/{vehicleId}
func (h *EventVehicleHandler) RemoveVehicle(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	vehicleIDStr := chi.URLParam(r, "vehicleId")
	vehicleID, err := uuid.Parse(vehicleIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	if err := h.vehicleService.RemoveVehicle(r.Context(), eventID, vehicleID); err != nil {
		if err == eventvehicle.ErrVehicleNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to remove vehicle from event")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Vehicle removed from event successfully")
}

// UpdateVehicleQuantity updates the quantity of a vehicle in an event
// PATCH /api/v1/events/{eventId}/vehicles/{vehicleId}/quantity
func (h *EventVehicleHandler) UpdateVehicleQuantity(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	vehicleIDStr := chi.URLParam(r, "vehicleId")
	vehicleID, err := uuid.Parse(vehicleIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	var req eventvehicle.UpdateVehicleQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.vehicleService.UpdateQuantity(r.Context(), eventID, vehicleID, &req); err != nil {
		switch err {
		case eventvehicle.ErrQuantityInvalid:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventvehicle.ErrVehicleNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to update vehicle quantity")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Vehicle quantity updated successfully")
}

// ListEventVehicles retrieves all vehicles for an event
// GET /api/v1/events/{eventId}/vehicles
func (h *EventVehicleHandler) ListEventVehicles(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	vehicles, err := h.vehicleService.ListByEvent(r.Context(), eventID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list vehicles")
		return
	}

	RespondSuccess(w, http.StatusOK, vehicles)
}
