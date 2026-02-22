package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/vehicle"
	"github.com/customermx/backend/internal/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// VehicleHandler handles vehicle endpoints
type VehicleHandler struct {
	vehicleService vehicle.Service
}

// NewVehicleHandler creates a new VehicleHandler
func NewVehicleHandler(vehicleService vehicle.Service) *VehicleHandler {
	return &VehicleHandler{vehicleService: vehicleService}
}

// CreateVehicle creates a new vehicle
// POST /api/v1/vehicles
func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req vehicle.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vehicleResp, err := h.vehicleService.CreateVehicle(r.Context(), &req)
	if err != nil {
		switch err {
		case vehicle.ErrVehicleAlreadyExists:
			RespondError(w, http.StatusConflict, err.Error())
		case vehicle.ErrBrandIDRequired, vehicle.ErrModelNameRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create vehicle")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, vehicleResp)
}

// GetVehicle retrieves a vehicle by ID
// GET /api/v1/vehicles/{id}
func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	vehicleResp, err := h.vehicleService.GetVehicle(r.Context(), id)
	if err != nil {
		if err == vehicle.ErrVehicleNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get vehicle")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, vehicleResp)
}

// UpdateVehicle updates an existing vehicle
// PUT /api/v1/vehicles/{id}
func (h *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	var req vehicle.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vehicleResp, err := h.vehicleService.UpdateVehicle(r.Context(), id, &req)
	if err != nil {
		switch err {
		case vehicle.ErrVehicleNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case vehicle.ErrModelNameRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to update vehicle")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, vehicleResp)
}

// DeleteVehicle deletes a vehicle
// DELETE /api/v1/vehicles/{id}
func (h *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	if err := h.vehicleService.DeleteVehicle(r.Context(), id); err != nil {
		if err == vehicle.ErrVehicleNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete vehicle")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Vehicle deleted successfully")
}

// ListVehicles retrieves all vehicles (BRAND users only see their brand's vehicles)
// GET /api/v1/vehicles
func (h *VehicleHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// BRAND users: only see vehicles from their brand
	if claims.Role == "BRAND" && claims.BrandID != nil {
		vehicles, err := h.vehicleService.ListVehiclesByBrand(r.Context(), *claims.BrandID)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to list vehicles")
			return
		}
		RespondSuccess(w, http.StatusOK, vehicles)
		return
	}

	vehicles, err := h.vehicleService.ListVehicles(r.Context())
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list vehicles")
		return
	}

	RespondSuccess(w, http.StatusOK, vehicles)
}

// ListVehiclesByBrand retrieves all vehicles for a specific brand
// GET /api/v1/brands/{brandId}/vehicles
func (h *VehicleHandler) ListVehiclesByBrand(w http.ResponseWriter, r *http.Request) {
	brandIDStr := chi.URLParam(r, "brandId")
	brandID, err := uuid.Parse(brandIDStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	vehicles, err := h.vehicleService.ListVehiclesByBrand(r.Context(), brandID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list vehicles by brand")
		return
	}

	RespondSuccess(w, http.StatusOK, vehicles)
}
