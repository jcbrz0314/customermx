package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/eventcoordinator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EventCoordinatorHandler handles event coordinator endpoints
type EventCoordinatorHandler struct {
	coordinatorService eventcoordinator.Service
}

// NewEventCoordinatorHandler creates a new EventCoordinatorHandler
func NewEventCoordinatorHandler(coordinatorService eventcoordinator.Service) *EventCoordinatorHandler {
	return &EventCoordinatorHandler{coordinatorService: coordinatorService}
}

// AssignCoordinator assigns a coordinator to an event
// POST /api/v1/events/{eventId}/coordinators
func (h *EventCoordinatorHandler) AssignCoordinator(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventId")

	var req eventcoordinator.AssignCoordinatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Override eventId from URL if present
	if eventID != "" {
		req.EventID = eventID
	}

	ec, err := h.coordinatorService.Assign(r.Context(), &req)
	if err != nil {
		switch err {
		case eventcoordinator.ErrEventIDRequired, eventcoordinator.ErrUserIDRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventcoordinator.ErrEventNotFound, eventcoordinator.ErrUserNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case eventcoordinator.ErrUserNotCoordinator:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventcoordinator.ErrCoordinatorAlreadyExists:
			RespondError(w, http.StatusConflict, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to assign coordinator")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, ec)
}

// RemoveCoordinator removes a coordinator from an event
// DELETE /api/v1/events/{eventId}/coordinators/{userId}
func (h *EventCoordinatorHandler) RemoveCoordinator(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.coordinatorService.Remove(r.Context(), eventID, userID); err != nil {
		if err == eventcoordinator.ErrCoordinatorNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to remove coordinator")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Coordinator removed successfully")
}

// ListEventCoordinators retrieves all coordinators for an event
// GET /api/v1/events/{eventId}/coordinators
func (h *EventCoordinatorHandler) ListEventCoordinators(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	coordinators, err := h.coordinatorService.ListByEvent(r.Context(), eventID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list coordinators")
		return
	}

	RespondSuccess(w, http.StatusOK, coordinators)
}

// ListCoordinatorEvents retrieves all events assigned to a coordinator
// GET /api/v1/coordinators/{userId}/events
func (h *EventCoordinatorHandler) ListCoordinatorEvents(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	events, err := h.coordinatorService.ListByUser(r.Context(), userID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}

	RespondSuccess(w, http.StatusOK, events)
}
