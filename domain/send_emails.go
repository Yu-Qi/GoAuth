package domain

// SendEmailService is an interface for sending emails
type SendEmailService interface {
	SendEmail(email string, subject string, body string) error
}
