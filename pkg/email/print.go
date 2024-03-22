package email

import "fmt"

// PrintEmailService is a service that prints the email to the console
type PrintEmailService struct{}

func NewPrintEmailService() *PrintEmailService {
	return &PrintEmailService{}
}

// SendEmail prints the email to the console
func (p PrintEmailService) SendEmail(email string, subject string, body string) error {
	fmt.Printf("Email: %s\nSubject: %s\nBody: %s\n", email, subject, body)
	return nil
}
