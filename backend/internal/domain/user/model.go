package user

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleCoordinator Role = "COORDINATOR"
	RoleBrand       Role = "BRAND"
	RoleVisualizer  Role = "VISUALIZER"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose in JSON
	Role         Role       `json:"role"`
	BrandID      *uuid.UUID `json:"brand_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Role     Role       `json:"role"`
	BrandID  *uuid.UUID `json:"brand_id,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Role     *Role   `json:"role,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

// UserResponse represents a user without sensitive data
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	BrandID   *uuid.UUID `json:"brand_id,omitempty"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts a User to UserResponse (removes sensitive data)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		BrandID:   u.BrandID,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.Role == "" {
		return ErrRoleRequired
	}

	// Validate role
	switch r.Role {
	case RoleAdmin, RoleCoordinator, RoleBrand, RoleVisualizer:
		// valid
	default:
		return ErrRoleRequired
	}

	// BRAND users must have a brand_id
	if r.Role == RoleBrand && r.BrandID == nil {
		return ErrBrandIDRequired
	}

	// Non-BRAND users must not have a brand_id
	if r.Role != RoleBrand && r.BrandID != nil {
		return ErrBrandIDNotAllowed
	}

	return nil
}
