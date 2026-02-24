package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// SMTPService implements the mail Service interface using SMTP
type SMTPService struct {
	config *Config
	host   string
	port   int
	auth   smtp.Auth
	from   string
}

// NewSMTPService creates a new SMTP email service
func NewSMTPService(config *Config) (*SMTPService, error) {
	if config.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if config.SMTPPort == 0 {
		return nil, fmt.Errorf("SMTP port is required")
	}
	if config.SMTPUsername == "" {
		return nil, fmt.Errorf("SMTP username is required")
	}
	if config.SMTPPassword == "" {
		return nil, fmt.Errorf("SMTP password is required")
	}

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	return &SMTPService{
		config: config,
		host:   config.SMTPHost,
		port:   config.SMTPPort,
		auth:   auth,
		from:   config.FromAddress,
	}, nil
}

// sendEmail sends a raw HTML email via SMTP with STARTTLS
func (s *SMTPService) sendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	tlsConfig := &tls.Config{ServerName: s.host}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	if err := client.Auth(s.auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	fromAddr := s.extractEmail(s.from)
	if err := client.Mail(fromAddr); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	msg := s.composeMessage(fromAddr, to, subject, body)
	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write email: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}

func (s *SMTPService) composeMessage(from, to, subject, body string) string {
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body
	return msg
}

func (s *SMTPService) extractEmail(addr string) string {
	if strings.Contains(addr, "<") && strings.Contains(addr, ">") {
		start := strings.Index(addr, "<") + 1
		end := strings.Index(addr, ">")
		return addr[start:end]
	}
	return addr
}

// SendInvitation sends an invitation email to a new user
func (s *SMTPService) SendInvitation(to, token, role string) error {
	inviteURL := fmt.Sprintf("%s/invite/accept?token=%s", s.config.FrontendURL, token)

	roleLabel := map[string]string{
		"ADMIN":       "Administrador",
		"COORDINATOR": "Coordinador",
		"BRAND":       "Marca",
	}
	label := roleLabel[role]
	if label == "" {
		label = role
	}

	content := paragraph("Has sido invitado a unirte a <strong>CustomerMX</strong>.") +
		paragraph(fmt.Sprintf("Tu rol en la plataforma será: <strong>%s</strong>.", label)) +
		paragraph("Haz clic en el botón para aceptar la invitación y crear tu contraseña.") +
		smallText("Este enlace expirará en 48 horas.")

	body := emailLayout(
		"Invitación a CustomerMX",
		content,
		"Aceptar Invitación",
		inviteURL,
		"#1976d2",
		s.config.LogoURL,
	)

	return s.sendEmail(to, "Invitación a CustomerMX", body)
}

// SendEventAssignment notifies a coordinator they were assigned to an event
func (s *SMTPService) SendEventAssignment(to string, details EventAssignmentDetails) error {
	eventURL := fmt.Sprintf("%s/events/%s", s.config.FrontendURL, details.EventID)

	durationLabel := "día"
	if details.DurationDays != 1 {
		durationLabel = "días"
	}

	rows := infoRow("Evento:", details.EventName) +
		infoRow("Marca:", details.BrandName) +
		infoRow("Tipo:", details.EventType) +
		infoRow("Organizador:", details.Organizer) +
		infoRow("Fecha de inicio:", details.StartDate) +
		infoRow("Duración:", fmt.Sprintf("%d %s", details.DurationDays, durationLabel)) +
		infoRow("Ubicación:", fmt.Sprintf("%s, %s", details.City, details.State)) +
		infoRow("Recinto:", details.Venue) +
		infoRow("Concesionario:", details.Dealer)

	content := paragraph("Has sido asignado como coordinador a un nuevo evento.") +
		infoBox(rows) +
		paragraph("Revisa los detalles e ingresa a la plataforma para comenzar la coordinación.")

	body := emailLayout(
		"Nueva Asignación de Evento",
		content,
		"Ver Evento",
		eventURL,
		"#1976d2",
		s.config.LogoURL,
	)

	return s.sendEmail(to, fmt.Sprintf("Asignado al evento: %s", details.EventName), body)
}

// SendEventCompleted notifies relevant users that an event has been completed
func (s *SMTPService) SendEventCompleted(to, eventName string) error {
	content := paragraph(fmt.Sprintf("El evento <strong>%s</strong> ha sido completado.", eventName)) +
		paragraph("La información del evento ya está disponible en la plataforma.")

	body := emailLayout(
		"Evento Completado",
		content,
		"Ver Reporte",
		fmt.Sprintf("%s/events", s.config.FrontendURL),
		"#2e7d32",
		s.config.LogoURL,
	)

	return s.sendEmail(to, fmt.Sprintf("Evento completado: %s", eventName), body)
}
