package notification

import "log"

type Service interface {
	Send(to, subject, body string) error
}

type EmailService struct {
	// smtp config would go here
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) Send(to, subject, body string) error {
	// Simulate sending email
	log.Printf("📧 Sending Email to %s | Subject: %s | Body: %s", to, subject, body)
	return nil
}
