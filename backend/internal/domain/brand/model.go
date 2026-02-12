package brand

import (
	"time"

	"github.com/google/uuid"
)

// Brand represents an automotive brand
type Brand struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateBrandRequest represents the request to create a new brand
type CreateBrandRequest struct {
	Name string `json:"name"`
}

// UpdateBrandRequest represents the request to update a brand
type UpdateBrandRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateBrandRequest
func (r *CreateBrandRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	return nil
}

// Validate validates the UpdateBrandRequest
func (r *UpdateBrandRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	return nil
}
