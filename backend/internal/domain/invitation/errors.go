package invitation

import "errors"

var (
	// Validation errors
	ErrEmailRequired        = errors.New("email is required")
	ErrRoleRequired         = errors.New("role is required")
	ErrBrandIDRequired      = errors.New("brand_id is required for BRAND role invitations")
	ErrBrandIDNotAllowed    = errors.New("brand_id is not allowed for non-BRAND invitations")
	ErrTokenRequired        = errors.New("token is required")
	ErrNameRequired         = errors.New("name is required")
	ErrPasswordRequired     = errors.New("password is required")

	// Business logic errors
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrInvitationExpired    = errors.New("invitation has expired")
	ErrInvitationAccepted   = errors.New("invitation already accepted")
	ErrUserAlreadyExists    = errors.New("user with this email already exists")
)
