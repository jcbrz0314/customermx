package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/user"
	"github.com/customermx/backend/internal/http/middleware"
	"github.com/customermx/backend/internal/infra/security"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userService user.Service
	jwtService  *security.JWTService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService user.Service, jwtService *security.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		switch err {
		case user.ErrInvalidCredentials:
			RespondError(w, http.StatusUnauthorized, err.Error())
		case user.ErrUserInactive:
			RespondError(w, http.StatusForbidden, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, response)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Get user to verify they still exist and are active
	userResp, err := h.userService.GetUser(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "User not found")
		return
	}

	if !userResp.IsActive {
		RespondError(w, http.StatusForbidden, "User account is inactive")
		return
	}

	// Generate new token pair
	tokens, err := h.jwtService.GenerateTokenPair(
		userResp.ID,
		string(userResp.Role),
		userResp.BrandID,
		userResp.IsActive,
	)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	RespondSuccess(w, http.StatusOK, map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT setup, logout is typically handled client-side
	// by deleting the tokens. For server-side logout, you would need to
	// implement token blacklisting or use a token version in the database.
	RespondSuccessWithMessage(w, http.StatusOK, nil, "Logged out successfully")
}

// GetCurrentUser returns the current authenticated user
// GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userResp, err := h.userService.GetUser(r.Context(), claims.UserID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	RespondSuccess(w, http.StatusOK, userResp)
}
