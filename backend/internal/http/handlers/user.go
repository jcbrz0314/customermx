package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService user.Service
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser creates a new user
// POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userResp, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		switch err {
		case user.ErrEmailAlreadyExists:
			RespondError(w, http.StatusConflict, err.Error())
		case user.ErrNameRequired, user.ErrEmailRequired, user.ErrPasswordRequired,
			user.ErrRoleRequired, user.ErrBrandIDRequired, user.ErrBrandIDNotAllowed:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create user")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, userResp)
}

// GetUser retrieves a user by ID
// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userResp, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get user")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, userResp)
}

// UpdateUser updates an existing user
// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userResp, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		if err == user.ErrUserNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, userResp)
}

// DeleteUser deletes a user
// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		if err == user.ErrUserNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "User deleted successfully")
}

// ListUsers retrieves all users
// GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	RespondSuccess(w, http.StatusOK, users)
}

// ListUsersByRole retrieves all users by role
// GET /api/v1/users/role/{role}
func (h *UserHandler) ListUsersByRole(w http.ResponseWriter, r *http.Request) {
	roleStr := chi.URLParam(r, "role")
	role := user.Role(roleStr)

	users, err := h.userService.ListUsersByRole(r.Context(), role)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list users by role")
		return
	}

	RespondSuccess(w, http.StatusOK, users)
}
