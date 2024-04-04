package email

import "github.com/Yu-Qi/GoAuth/domain"

var (
	service domain.SendEmailService
)

// GetService returns the email service
func GetService() domain.SendEmailService {
	return service
}

// InitService initializes the email service
func InitService(s domain.SendEmailService) {
	service = s
}
