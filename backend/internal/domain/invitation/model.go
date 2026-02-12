package invitation

import (
	"time"

	"github.com/customermx/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Invitation represents a user invitation
type Invitation struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Role      user.Role   `json:"role"`
	BrandID   *uuid.UUID  `json:"brand_id,omitempty"`
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	Accepted  bool        `json:"accepted"`
	CreatedAt time.Time   `json:"created_at"`
}

// CreateInvitationRequest represents the request to create a new invitation
type CreateInvitationRequest struct {
	Email   string      `json:"email"`
	Role    user.Role   `json:"role"`
	BrandID *uuid.UUID  `json:"brand_id,omitempty"`
}

// AcceptInvitationRequest represents the request to accept an invitation
type AcceptInvitationRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Validate validates the CreateInvitationRequest
func (r *CreateInvitationRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Role == "" {
		return ErrRoleRequired
	}

	// BRAND invitations must have a brand_id
	if r.Role == user.RoleBrand && r.BrandID == nil {
		return ErrBrandIDRequired
	}

	// Non-BRAND invitations must not have a brand_id
	if r.Role != user.RoleBrand && r.BrandID != nil {
		return ErrBrandIDNotAllowed
	}

	return nil
}

// Validate validates the AcceptInvitationRequest
func (r *AcceptInvitationRequest) Validate() error {
	if r.Token == "" {
		return ErrTokenRequired
	}
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
