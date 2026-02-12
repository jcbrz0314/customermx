package mail

import (
	"fmt"
)

// Service defines the email service interface
type Service interface {
	SendInvitation(to, token, role string) error
	SendEventAssignment(to, eventName string) error
	SendEventCompleted(to, eventName string) error
}

// Config holds email service configuration
type Config struct {
	Provider     string
	FromAddress  string
	FrontendURL  string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPUseTLS   bool
}

// MockService is a mock email service for development
type MockService struct {
	config *Config
}

// NewMockService creates a new mock email service
func NewMockService(config *Config) *MockService {
	return &MockService{config: config}
}

// SendInvitation sends an invitation email (mock)
func (s *MockService) SendInvitation(to, token, role string) error {
	inviteURL := fmt.Sprintf("%s/invite/accept?token=%s", s.config.FrontendURL, token)
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("From: %s\n", s.config.FromAddress)
	fmt.Printf("Subject: Invitación a CustomerMX\n")
	fmt.Printf("Body:\n")
	fmt.Printf("Has sido invitado a CustomerMX como %s.\n", role)
	fmt.Printf("Haz clic en el siguiente enlace para aceptar:\n")
	fmt.Printf("%s\n", inviteURL)
	fmt.Printf("------------------\n\n")
	return nil
}

// SendEventAssignment sends an event assignment notification (mock)
func (s *MockService) SendEventAssignment(to, eventName string) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("From: %s\n", s.config.FromAddress)
	fmt.Printf("Subject: Asignación a evento\n")
	fmt.Printf("Body:\n")
	fmt.Printf("Has sido asignado al evento: %s\n", eventName)
	fmt.Printf("------------------\n\n")
	return nil
}

// SendEventCompleted sends an event completion notification (mock)
func (s *MockService) SendEventCompleted(to, eventName string) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\n", to)
	fmt.Printf("From: %s\n", s.config.FromAddress)
	fmt.Printf("Subject: Evento completado\n")
	fmt.Printf("Body:\n")
	fmt.Printf("El evento %s ha sido completado y la información está disponible.\n", eventName)
	fmt.Printf("------------------\n\n")
	return nil
}
