package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/invitation"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// InvitationHandler handles invitation endpoints
type InvitationHandler struct {
	invitationService invitation.Service
}

// NewInvitationHandler creates a new InvitationHandler
func NewInvitationHandler(invitationService invitation.Service) *InvitationHandler {
	return &InvitationHandler{invitationService: invitationService}
}

// CreateInvitation creates a new invitation
// POST /api/v1/invitations
func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	var req invitation.CreateInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	inv, err := h.invitationService.CreateInvitation(r.Context(), &req)
	if err != nil {
		switch err {
		case invitation.ErrUserAlreadyExists:
			RespondError(w, http.StatusConflict, err.Error())
		case invitation.ErrEmailRequired, invitation.ErrRoleRequired,
			invitation.ErrBrandIDRequired, invitation.ErrBrandIDNotAllowed:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create invitation")
		}
		return
	}

	// Don't expose the token in the response (it's sent via email)
	response := map[string]interface{}{
		"id":         inv.ID,
		"email":      inv.Email,
		"role":       inv.Role,
		"brand_id":   inv.BrandID,
		"expires_at": inv.ExpiresAt,
		"created_at": inv.CreatedAt,
	}

	RespondSuccessWithMessage(w, http.StatusCreated, response, "Invitation sent successfully")
}

// AcceptInvitation accepts an invitation and creates a user
// POST /api/v1/invitations/accept
func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req invitation.AcceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.invitationService.AcceptInvitation(r.Context(), &req)
	if err != nil {
		switch err {
		case invitation.ErrInvitationNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case invitation.ErrInvitationExpired:
			RespondError(w, http.StatusGone, err.Error())
		case invitation.ErrInvitationAccepted:
			RespondError(w, http.StatusConflict, err.Error())
		case invitation.ErrTokenRequired, invitation.ErrNameRequired, invitation.ErrPasswordRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to accept invitation")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusCreated, user.ToResponse(), "Account created successfully")
}

// GetInvitation retrieves an invitation by ID
// GET /api/v1/invitations/{id}
func (h *InvitationHandler) GetInvitation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	inv, err := h.invitationService.GetInvitation(r.Context(), id)
	if err != nil {
		if err == invitation.ErrInvitationNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get invitation")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, inv)
}

// ListInvitations retrieves all invitations
// GET /api/v1/invitations
func (h *InvitationHandler) ListInvitations(w http.ResponseWriter, r *http.Request) {
	// Check query parameter for pending only
	pendingOnly := r.URL.Query().Get("pending") == "true"

	var invitations []*invitation.Invitation
	var err error

	if pendingOnly {
		invitations, err = h.invitationService.ListPendingInvitations(r.Context())
	} else {
		invitations, err = h.invitationService.ListInvitations(r.Context())
	}

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list invitations")
		return
	}

	RespondSuccess(w, http.StatusOK, invitations)
}

// ResendInvitation resends an invitation email
// POST /api/v1/invitations/{id}/resend
func (h *InvitationHandler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	if err := h.invitationService.ResendInvitation(r.Context(), id); err != nil {
		switch err {
		case invitation.ErrInvitationNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case invitation.ErrInvitationAccepted:
			RespondError(w, http.StatusConflict, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to resend invitation")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Invitation resent successfully")
}

// ValidateInvitationToken validates an invitation token without authentication
// GET /api/v1/invitations/validate?token=xxx
func (h *InvitationHandler) ValidateInvitationToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		RespondError(w, http.StatusBadRequest, "Token is required")
		return
	}

	inv, err := h.invitationService.GetInvitationByToken(r.Context(), token)
	if err != nil {
		switch err {
		case invitation.ErrInvitationNotFound:
			RespondError(w, http.StatusNotFound, "Invitación no encontrada o inválida")
		case invitation.ErrInvitationExpired:
			RespondError(w, http.StatusGone, "La invitación ha expirado")
		case invitation.ErrInvitationAccepted:
			RespondError(w, http.StatusConflict, "La invitación ya fue aceptada")
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to validate token")
		}
		return
	}

	response := map[string]interface{}{
		"email":      inv.Email,
		"role":       inv.Role,
		"brand_id":   inv.BrandID,
		"expires_at": inv.ExpiresAt,
	}

	RespondSuccess(w, http.StatusOK, response)
}

// DeleteInvitation deletes an invitation
// DELETE /api/v1/invitations/{id}
func (h *InvitationHandler) DeleteInvitation(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid invitation ID")
		return
	}

	if err := h.invitationService.DeleteInvitation(r.Context(), id); err != nil {
		if err == invitation.ErrInvitationNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete invitation")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Invitation deleted successfully")
}
