package mail

import (
	"fmt"
	"log"
)

// EventAssignmentDetails holds the event info sent in assignment emails
type EventAssignmentDetails struct {
	EventID      string
	EventName    string
	BrandName    string
	EventType    string
	Organizer    string
	StartDate    string
	DurationDays int
	State        string
	City         string
	Venue        string
	Dealer       string
}

// EventUnassignmentDetails holds the event info sent in unassignment emails
type EventUnassignmentDetails struct {
	EventName    string
	BrandName    string
	StartDate    string
	City         string
	State        string
}

// Service defines the email service interface
type Service interface {
	SendInvitation(to, token, role string) error
	SendEventAssignment(to string, details EventAssignmentDetails) error
	SendEventUnassignment(to string, details EventUnassignmentDetails) error
	SendEventCompleted(to, eventName string) error
}

// Config holds email service configuration
type Config struct {
	Provider     string
	FromAddress  string
	FrontendURL  string
	LogoURL      string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPUseTLS   bool
}

// NewService creates the appropriate mail service based on provider config.
// Falls back to MockService if SMTP is configured but fails to connect.
func NewService(config *Config) Service {
	if config.Provider == "smtp" {
		svc, err := NewSMTPService(config)
		if err != nil {
			log.Printf("[mail] SMTP setup failed (%v), falling back to mock", err)
			return NewMockService(config)
		}
		log.Printf("[mail] SMTP service ready (%s:%d)", config.SMTPHost, config.SMTPPort)
		return svc
	}
	log.Printf("[mail] Using mock mail service")
	return NewMockService(config)
}

// MockService is a mock email service for development
type MockService struct {
	config *Config
}

// NewMockService creates a new mock email service
func NewMockService(config *Config) *MockService {
	return &MockService{config: config}
}

func (s *MockService) SendInvitation(to, token, role string) error {
	inviteURL := fmt.Sprintf("%s/invite/accept?token=%s", s.config.FrontendURL, token)
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\nSubject: Invitación a CustomerMX\n", to)
	fmt.Printf("Role: %s\nLink: %s\n", role, inviteURL)
	fmt.Printf("------------------\n\n")
	return nil
}

func (s *MockService) SendEventAssignment(to string, details EventAssignmentDetails) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\nSubject: Asignado al evento: %s\n", to, details.EventName)
	fmt.Printf("Marca: %s | Tipo: %s | Fecha: %s | %s, %s\n",
		details.BrandName, details.EventType, details.StartDate, details.City, details.State)
	fmt.Printf("------------------\n\n")
	return nil
}

func (s *MockService) SendEventUnassignment(to string, details EventUnassignmentDetails) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\nSubject: Removido del evento: %s\n", to, details.EventName)
	fmt.Printf("Marca: %s | Fecha: %s | %s, %s\n",
		details.BrandName, details.StartDate, details.City, details.State)
	fmt.Printf("------------------\n\n")
	return nil
}

func (s *MockService) SendEventCompleted(to, eventName string) error {
	fmt.Printf("\n--- MOCK EMAIL ---\n")
	fmt.Printf("To: %s\nSubject: Evento completado: %s\n", to, eventName)
	fmt.Printf("------------------\n\n")
	return nil
}
