package router

import (
	"net/http"
	"time"

	"github.com/customermx/backend/internal/config"
	"github.com/customermx/backend/internal/domain/brand"
	"github.com/customermx/backend/internal/domain/invitation"
	"github.com/customermx/backend/internal/domain/user"
	"github.com/customermx/backend/internal/domain/vehicle"
	"github.com/customermx/backend/internal/http/handlers"
	"github.com/customermx/backend/internal/http/middleware"
	"github.com/customermx/backend/internal/infra/db"
	"github.com/customermx/backend/internal/infra/mail"
	"github.com/customermx/backend/internal/infra/security"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// New creates a new HTTP router with all routes and middleware
func New(cfg *config.Config, dbConn *db.Connection) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize security services
	jwtService := security.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	passwordService := security.NewPasswordService()

	// Initialize mail service
	mailConfig := &mail.Config{
		Provider:    cfg.Email.Provider,
		FromAddress: cfg.Email.From,
		FrontendURL: cfg.Email.FrontendURL,
	}
	mailService := mail.NewMockService(mailConfig)

	// Initialize repositories
	userRepo := user.NewRepository(dbConn.Pool)
	brandRepo := brand.NewRepository(dbConn.Pool)
	vehicleRepo := vehicle.NewRepository(dbConn.Pool)
	invitationRepo := invitation.NewRepository(dbConn.Pool)

	// Initialize domain services
	userService := user.NewService(userRepo, jwtService, passwordService)
	brandService := brand.NewService(brandRepo)
	vehicleService := vehicle.NewService(vehicleRepo)
	invitationService := invitation.NewService(invitationRepo, userService, mailService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtService)
	userHandler := handlers.NewUserHandler(userService)
	brandHandler := handlers.NewBrandHandler(brandService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"customermx-api"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Group(func(r chi.Router) {
			// Auth routes
			r.Post("/auth/login", authHandler.Login)

			// Invitation routes (public)
			r.Post("/invitations/accept", invitationHandler.AcceptInvitation)
		})

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg))

			// Auth routes (protected)
			r.Post("/auth/refresh", authHandler.RefreshToken)
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/me", authHandler.GetCurrentUser)

			// User routes
			r.Get("/users", userHandler.ListUsers)
			r.Post("/users", userHandler.CreateUser)
			r.Get("/users/{id}", userHandler.GetUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Delete("/users/{id}", userHandler.DeleteUser)
			r.Get("/users/role/{role}", userHandler.ListUsersByRole)

			// Brand routes
			r.Get("/brands", brandHandler.ListBrands)
			r.Post("/brands", brandHandler.CreateBrand)
			r.Get("/brands/{id}", brandHandler.GetBrand)
			r.Put("/brands/{id}", brandHandler.UpdateBrand)
			r.Delete("/brands/{id}", brandHandler.DeleteBrand)

			// Vehicle routes
			r.Get("/vehicles", vehicleHandler.ListVehicles)
			r.Post("/vehicles", vehicleHandler.CreateVehicle)
			r.Get("/vehicles/{id}", vehicleHandler.GetVehicle)
			r.Put("/vehicles/{id}", vehicleHandler.UpdateVehicle)
			r.Delete("/vehicles/{id}", vehicleHandler.DeleteVehicle)
			r.Get("/brands/{brandId}/vehicles", vehicleHandler.ListVehiclesByBrand)

			// Invitation routes (protected)
			r.Get("/invitations", invitationHandler.ListInvitations)
			r.Post("/invitations", invitationHandler.CreateInvitation)
			r.Get("/invitations/{id}", invitationHandler.GetInvitation)
			r.Post("/invitations/{id}/resend", invitationHandler.ResendInvitation)
			r.Delete("/invitations/{id}", invitationHandler.DeleteInvitation)
		})
	})

	return r
}
