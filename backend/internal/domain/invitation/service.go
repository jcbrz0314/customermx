package invitation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/customermx/backend/internal/domain/user"
	"github.com/customermx/backend/internal/infra/mail"
	"github.com/google/uuid"
)

// Service defines the invitation business logic interface
type Service interface {
	CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*Invitation, error)
	GetInvitation(ctx context.Context, id uuid.UUID) (*Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*Invitation, error)
	AcceptInvitation(ctx context.Context, req *AcceptInvitationRequest) (*user.User, error)
	ListInvitations(ctx context.Context) ([]*Invitation, error)
	ListPendingInvitations(ctx context.Context) ([]*Invitation, error)
	ResendInvitation(ctx context.Context, id uuid.UUID) error
	DeleteInvitation(ctx context.Context, id uuid.UUID) error
}

// InvitationService implements the Service interface
type InvitationService struct {
	repo        Repository
	userService user.Service
	mailService mail.Service
}

// NewService creates a new InvitationService
func NewService(repo Repository, userService user.Service, mailService mail.Service) *InvitationService {
	return &InvitationService{
		repo:        repo,
		userService: userService,
		mailService: mailService,
	}
}

// CreateInvitation creates a new invitation and sends email
func (s *InvitationService) CreateInvitation(ctx context.Context, req *CreateInvitationRequest) (*Invitation, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Generate secure token
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// Create invitation
	invitation := &Invitation{
		ID:        uuid.New(),
		Email:     req.Email,
		Role:      req.Role,
		BrandID:   req.BrandID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		Accepted:  false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	// Send invitation email
	if err := s.mailService.SendInvitation(req.Email, token, string(req.Role)); err != nil {
		// Log error but don't fail - invitation is created
		// TODO: Add proper logging
	}

	return invitation, nil
}

// GetInvitation retrieves an invitation by ID
func (s *InvitationService) GetInvitation(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	return s.repo.GetByID(ctx, id)
}

// GetInvitationByToken retrieves an invitation by token
func (s *InvitationService) GetInvitationByToken(ctx context.Context, token string) (*Invitation, error) {
	invitation, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if time.Now().After(invitation.ExpiresAt) {
		return nil, ErrInvitationExpired
	}

	// Check if already accepted
	if invitation.Accepted {
		return nil, ErrInvitationAccepted
	}

	return invitation, nil
}

// AcceptInvitation accepts an invitation and creates a user
func (s *InvitationService) AcceptInvitation(ctx context.Context, req *AcceptInvitationRequest) (*user.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get invitation
	invitation, err := s.GetInvitationByToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	// Create user
	createUserReq := &user.CreateUserRequest{
		Name:     req.Name,
		Email:    invitation.Email,
		Password: req.Password,
		Role:     invitation.Role,
		BrandID:  invitation.BrandID,
	}

	userResp, err := s.userService.CreateUser(ctx, createUserReq)
	if err != nil {
		return nil, err
	}

	// Mark invitation as accepted
	invitation.Accepted = true
	if err := s.repo.Update(ctx, invitation); err != nil {
		// Log error but don't fail - user is already created
		// TODO: Add proper logging
	}

	// Convert UserResponse back to User (we need the full user for response)
	// In a real scenario, you might want to adjust the service to return the full user
	fullUser := &user.User{
		ID:        userResp.ID,
		Name:      userResp.Name,
		Email:     userResp.Email,
		Role:      userResp.Role,
		BrandID:   userResp.BrandID,
		IsActive:  userResp.IsActive,
		CreatedAt: userResp.CreatedAt,
		UpdatedAt: userResp.UpdatedAt,
	}

	return fullUser, nil
}

// ListInvitations retrieves all invitations
func (s *InvitationService) ListInvitations(ctx context.Context) ([]*Invitation, error) {
	return s.repo.List(ctx)
}

// ListPendingInvitations retrieves all pending invitations
func (s *InvitationService) ListPendingInvitations(ctx context.Context) ([]*Invitation, error) {
	return s.repo.ListPending(ctx)
}

// ResendInvitation resends an invitation email
func (s *InvitationService) ResendInvitation(ctx context.Context, id uuid.UUID) error {
	invitation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if invitation.Accepted {
		return ErrInvitationAccepted
	}

	// Extend expiration
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)
	if err := s.repo.Update(ctx, invitation); err != nil {
		return err
	}

	// Resend email
	return s.mailService.SendInvitation(invitation.Email, invitation.Token, string(invitation.Role))
}

// DeleteInvitation deletes an invitation
func (s *InvitationService) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// generateToken generates a secure random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
