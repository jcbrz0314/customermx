package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/eventreport"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// EventReportHandler handles event report endpoints
type EventReportHandler struct {
	reportService eventreport.Service
}

// NewEventReportHandler creates a new EventReportHandler
func NewEventReportHandler(reportService eventreport.Service) *EventReportHandler {
	return &EventReportHandler{reportService: reportService}
}

// CreateOrUpdateReport creates or updates an event report (UPSERT)
// POST /api/v1/events/{eventId}/report
func (h *EventReportHandler) CreateOrUpdateReport(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req eventreport.CreateOrUpdateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	report, err := h.reportService.CreateOrUpdate(r.Context(), eventID, &req)
	if err != nil {
		switch err {
		case eventreport.ErrDealerRatingInvalid:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventreport.ErrEventNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create or update report")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, report)
}

// GetEventReport retrieves an event report by event ID
// GET /api/v1/events/{eventId}/report
func (h *EventReportHandler) GetEventReport(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	report, err := h.reportService.GetByEventID(r.Context(), eventID)
	if err != nil {
		if err == eventreport.ErrReportNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get report")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, report)
}

// CompleteReport marks a report as completed or not completed
// PATCH /api/v1/events/{eventId}/report/complete
func (h *EventReportHandler) CompleteReport(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req eventreport.CompleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	report, err := h.reportService.MarkAsCompleted(r.Context(), eventID, &req)
	if err != nil {
		switch err {
		case eventreport.ErrReportNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case eventreport.ErrEventNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case eventreport.ErrEventNotCompleted:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to complete report")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, report)
}

// DeleteReport deletes an event report
// DELETE /api/v1/events/{eventId}/report
func (h *EventReportHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "eventId")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := h.reportService.Delete(r.Context(), eventID); err != nil {
		if err == eventreport.ErrReportNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete report")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Report deleted successfully")
}
