package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/customermx/backend/internal/domain/brand"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// BrandHandler handles brand endpoints
type BrandHandler struct {
	brandService brand.Service
}

// NewBrandHandler creates a new BrandHandler
func NewBrandHandler(brandService brand.Service) *BrandHandler {
	return &BrandHandler{brandService: brandService}
}

// CreateBrand creates a new brand
// POST /api/v1/brands
func (h *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var req brand.CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	brandResp, err := h.brandService.CreateBrand(r.Context(), &req)
	if err != nil {
		switch err {
		case brand.ErrBrandAlreadyExists:
			RespondError(w, http.StatusConflict, err.Error())
		case brand.ErrNameRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to create brand")
		}
		return
	}

	RespondSuccess(w, http.StatusCreated, brandResp)
}

// GetBrand retrieves a brand by ID
// GET /api/v1/brands/{id}
func (h *BrandHandler) GetBrand(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	brandResp, err := h.brandService.GetBrand(r.Context(), id)
	if err != nil {
		if err == brand.ErrBrandNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to get brand")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, brandResp)
}

// UpdateBrand updates an existing brand
// PUT /api/v1/brands/{id}
func (h *BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	var req brand.UpdateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	brandResp, err := h.brandService.UpdateBrand(r.Context(), id, &req)
	if err != nil {
		switch err {
		case brand.ErrBrandNotFound:
			RespondError(w, http.StatusNotFound, err.Error())
		case brand.ErrNameRequired:
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to update brand")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, brandResp)
}

// DeleteBrand deletes a brand
// DELETE /api/v1/brands/{id}
func (h *BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	if err := h.brandService.DeleteBrand(r.Context(), id); err != nil {
		if err == brand.ErrBrandNotFound {
			RespondError(w, http.StatusNotFound, err.Error())
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete brand")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Brand deleted successfully")
}

// ListBrands retrieves all brands
// GET /api/v1/brands
func (h *BrandHandler) ListBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.brandService.ListBrands(r.Context())
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list brands")
		return
	}

	RespondSuccess(w, http.StatusOK, brands)
}
