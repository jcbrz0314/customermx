package middleware

import (
	"net/http"

	"github.com/customermx/backend/internal/domain/event"
	"github.com/customermx/backend/internal/domain/eventcoordinator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// RequireRole is a middleware that restricts access to specific roles
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// VISUALIZER is read-only: block all write operations regardless of allowedRoles
			if claims.Role == "VISUALIZER" {
				http.Error(w, "Forbidden: VISUALIZER role cannot perform write operations", http.StatusForbidden)
				return
			}

			// Check if user's role is in the allowed roles
			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// RequireEventAccess is a middleware that checks if the user has access to a specific event
// ADMIN: full access
// COORDINATOR: only assigned events
// BRAND: only events of their brand
func RequireEventAccess(eventService event.Service, coordinatorService eventcoordinator.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ADMIN and VISUALIZER have full read access
			if claims.Role == "ADMIN" || claims.Role == "VISUALIZER" {
				next.ServeHTTP(w, r)
				return
			}

			// Get event ID from URL
			eventIDStr := chi.URLParam(r, "id")
			if eventIDStr == "" {
				eventIDStr = chi.URLParam(r, "eventId")
			}

			eventID, err := uuid.Parse(eventIDStr)
			if err != nil {
				http.Error(w, "Invalid event ID", http.StatusBadRequest)
				return
			}

			// Get event details
			evt, err := eventService.GetByID(r.Context(), eventID)
			if err != nil {
				if err == event.ErrEventNotFound {
					http.Error(w, "Event not found", http.StatusNotFound)
				} else {
					http.Error(w, "Failed to get event", http.StatusInternalServerError)
				}
				return
			}

			// COORDINATOR: check if assigned to this event
			if claims.Role == "COORDINATOR" {
				assigned, err := coordinatorService.IsCoordinatorAssigned(r.Context(), eventID, claims.UserID)
				if err != nil {
					http.Error(w, "Failed to check coordinator assignment", http.StatusInternalServerError)
					return
				}

				if assigned {
					next.ServeHTTP(w, r)
					return
				}

				http.Error(w, "Forbidden: not assigned to this event", http.StatusForbidden)
				return
			}

			// BRAND: check if event belongs to their brand
			if claims.Role == "BRAND" {
				if claims.BrandID != nil && evt.BrandID == *claims.BrandID {
					next.ServeHTTP(w, r)
					return
				}

				http.Error(w, "Forbidden: event does not belong to your brand", http.StatusForbidden)
				return
			}

			// Default: deny access
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}
