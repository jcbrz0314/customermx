package user

import (
	"context"
	"time"

	"github.com/customermx/backend/internal/infra/security"
	"github.com/google/uuid"
)

// Service defines the user business logic interface
type Service interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]*UserResponse, error)
	ListUsersByBrand(ctx context.Context, brandID uuid.UUID) ([]*UserResponse, error)
	ListUsersByRole(ctx context.Context, role Role) ([]*UserResponse, error)
	ActivateUser(ctx context.Context, email, password string) (*User, error)
}

// UserService implements the Service interface
type UserService struct {
	repo            Repository
	jwtService      *security.JWTService
	passwordService *security.PasswordService
}

// NewService creates a new UserService
func NewService(repo Repository, jwtService *security.JWTService, passwordService *security.PasswordService) *UserService {
	return &UserService{
		repo:            repo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

// Login authenticates a user and returns tokens
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := s.passwordService.Verify(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT tokens
	tokens, err := s.jwtService.GenerateTokenPair(user.ID, string(user.Role), user.BrandID, user.IsActive)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Hash password
	passwordHash, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		BrandID:      req.BrandID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ListUsers retrieves all users
func (s *UserService) ListUsers(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// ListUsersByBrand retrieves all users for a specific brand
func (s *UserService) ListUsersByBrand(ctx context.Context, brandID uuid.UUID) ([]*UserResponse, error) {
	users, err := s.repo.ListByBrandID(ctx, brandID)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// ListUsersByRole retrieves all users with a specific role
func (s *UserService) ListUsersByRole(ctx context.Context, role Role) ([]*UserResponse, error) {
	users, err := s.repo.ListByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// ActivateUser activates a user account (used when accepting invitation)
func (s *UserService) ActivateUser(ctx context.Context, email, password string) (*User, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Hash new password
	passwordHash, err := s.passwordService.Hash(password)
	if err != nil {
		return nil, err
	}

	// Update user with new password and activate
	user.PasswordHash = passwordHash
	user.IsActive = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
